package main

import (
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func run() {
	router := gin.New()
	router.TrustedPlatform = "X-Real-IP"
	server.Handler = router

	router.Use(gin.Recovery())

	if logPath != "" {
		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
		if err != nil {
			log.Fatalln("Failed to open log file:", err)
		}

		gin.DefaultWriter = f
		gin.DefaultErrorWriter = f

		log.SetOutput(f)

		router.Use(gin.Logger())
	} else {
		log.SetOutput(io.Discard)
	}

	if err := initDB(); err != nil {
		log.Fatalln("Failed to initialize database:", err)
	}

	js, err := os.ReadFile(joinPath(dir(self), "dist/const.js"))
	if err != nil {
		log.Fatal(err)
	}

	if *universal {
		var redisStore struct{ Endpoint, Password, Secret, API string }
		if err := meta.Get("account_redis", &redisStore); err != nil {
			log.Fatal(err)
		}

		js = bytes.ReplaceAll(js, []byte("@universal@"), []byte(redisStore.API))

		store, err := redis.NewStore(10, "tcp", redisStore.Endpoint, redisStore.Password, []byte(redisStore.Secret))
		if err != nil {
			log.Fatal(err)
		}
		if err := redis.SetKeyPrefix(store, "account_"); err != nil {
			log.Fatal(err)
		}
		router.Use(sessions.Sessions("universal", store))
	} else {
		js = bytes.ReplaceAll(js, []byte("@universal@"), nil)

		secret := make([]byte, 16)
		if _, err := rand.Read(secret); err != nil {
			log.Fatalln("Failed to get secret:", err)
		}
		router.Use(sessions.Sessions("session", cookie.NewStore(secret)))
	}

	if priv != nil {
		pubkey_bytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
		if err != nil {
			log.Fatal(err)
		}
		js = bytes.ReplaceAll(
			js, []byte("@pubkey@"),
			bytes.ReplaceAll(
				pem.EncodeToMemory(&pem.Block{
					Type:  "RSA PUBLIC KEY",
					Bytes: pubkey_bytes,
				}),
				[]byte{'\n'},
				nil,
			),
		)
	} else {
		js = bytes.ReplaceAll(js, []byte("@pubkey@"), nil)
	}

	if err := os.WriteFile(joinPath(dir(self), "dist/env.js"), js, 0644); err != nil {
		log.Fatal(err)
	}

	router.StaticFS("/assets", http.Dir(joinPath(dir(self), "dist/assets")))
	router.StaticFile("env.js", joinPath(dir(self), "dist/env.js"))
	router.StaticFile("favicon.ico", joinPath(dir(self), "dist/favicon.ico"))
	router.LoadHTMLFiles(joinPath(dir(self), "dist/index.html"))

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	router.GET("/info", func(c *gin.Context) {
		_, username, _ := getUser(c)
		c.JSON(200, gin.H{"username": username, "refresh": refresh})
	})

	if !*universal {
		auth := router.Group("/")
		auth.POST("/login", login)
		auth.POST("/logout", authRequired, func(c *gin.Context) {
			session := sessions.Default(c)
			session.Clear()
			session.Options(sessions.Options{MaxAge: -1})
			if err := session.Save(); err != nil {
				log.Print(err)
				c.String(500, "")
				return
			}
			c.String(200, "bye")
		})
		auth.POST("/chgpwd", authRequired, chgpwd)
	}

	base := router.Group("/")
	base.GET("/stock/:index/:code", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	base.GET("/mystocks", myStocks)
	base.GET("/flows", capitalFlows)
	base.GET("/indices", indices)
	base.POST("/realtime", getRealtime)
	base.POST("/chart", getChart)
	base.POST("/suggest", getSuggest)
	base.GET("/star", star)
	base.POST("/star", doStar)
	base.POST("/reorder", reorder)
	base.POST("/refresh", func(c *gin.Context) {
		userID, _, _ := getUser(c)
		if userID != "" {
			loadStocks(userID, true)
		}
		c.String(200, "")
	})

	router.NoRoute(func(c *gin.Context) {
		c.Redirect(302, "/")
	})

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
