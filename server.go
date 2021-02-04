package main

import (
	"crypto/rand"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func run() {
	if *logPath != "" {
		f, err := os.OpenFile(*logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
		if err != nil {
			log.Fatalln("Failed to open log file:", err)
		}
		gin.DefaultWriter = f
		gin.DefaultErrorWriter = f
		log.SetOutput(f)
	}

	if err := initDB(); err != nil {
		log.Fatalln("Failed to initialize database:", err)
	}

	router := gin.Default()
	server.Handler = router

	js, err := ioutil.ReadFile(joinPath(dir(self), "public/build/bundle.js"))
	if err != nil {
		log.Fatal(err)
	}

	if universal {
		var redisStore struct{ Endpoint, Password, Secret, API string }
		if err := meta.Get("account_redis", &redisStore); err != nil {
			log.Fatal(err)
		}

		if err := ioutil.WriteFile(joinPath(dir(self), "public/build/script.js"),
			[]byte(strings.ReplaceAll(string(js), "@universal@", redisStore.API)), 0644); err != nil {
			log.Fatal(err)
		}

		store, err := redis.NewStore(10, "tcp", redisStore.Endpoint, redisStore.Password, []byte(redisStore.Secret))
		if err != nil {
			log.Fatal(err)
		}
		router.Use(sessions.Sessions("session", store))
	} else {
		if err := ioutil.WriteFile(joinPath(dir(self), "public/build/script.js"),
			[]byte(strings.ReplaceAll(string(js), "@universal@", "")), 0644); err != nil {
			log.Fatal(err)
		}

		secret := make([]byte, 16)
		if _, err := rand.Read(secret); err != nil {
			log.Fatalln("Failed to get secret:", err)
		}
		router.Use(sessions.Sessions("session", cookie.NewStore(secret)))
	}

	router.StaticFS("/build", http.Dir(joinPath(dir(self), "public/build")))
	router.StaticFile("favicon.ico", joinPath(dir(self), "public/favicon.ico"))
	router.LoadHTMLFiles(joinPath(dir(self), "public/index.html"))

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	router.GET("/info", func(c *gin.Context) {
		_, username, _ := getUser(c)
		c.JSON(200, gin.H{"username": username, "refresh": refresh})
	})

	if !universal {
		auth := router.Group("/")
		auth.POST("/login", login)
		auth.GET("/logout", authRequired, func(c *gin.Context) {
			session := sessions.Default(c)
			session.Clear()
			if err := session.Save(); err != nil {
				log.Print(err)
				c.String(500, "")
				return
			}
			c.Redirect(302, "/")
		})
		auth.POST("/chgpwd", authRequired, chgpwd)
	}

	base := router.Group("/")
	base.GET("/stock/:index/:code", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	base.GET("/mystocks", myStocks)
	base.GET("/indices", indices)
	base.POST("/get", getStock)
	base.POST("/suggest", getSuggest)
	base.GET("/star", star)
	base.POST("/star", doStar)
	base.POST("/reorder", reorder)

	router.NoRoute(func(c *gin.Context) {
		c.Redirect(302, "/")
	})

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
