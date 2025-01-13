package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/database/mongodb"
	"github.com/sunshineplan/password"
)

type user struct {
	ID       mongodb.OID `json:"_id" bson:"_id"`
	Username string
	Password string
}

func authRequired(c *gin.Context) {
	if sessions.Default(c).Get("id") == nil {
		c.Abort()
		c.Redirect(302, "/")
	}
}

func getUser(c *gin.Context) (id, username string, err error) {
	session := sessions.Default(c)
	sid := session.Get("id")
	if sid == nil {
		return
	}
	username, _ = session.Get("username").(string)
	if *universal {
		var user user
		if err = accountClient.FindOne(mongodb.M{"uid": sid}, nil, &user); err != nil {
			return
		}
		id = user.ID.Hex()
		return
	}
	id, _ = sid.(string)
	return
}

type info struct {
	username string
	ip       string
}

func auth(id info, got, want string) (message string, err error) {
	if err = password.CompareHashAndPassword(id, got, want); err != nil {
		if errors.Is(err, password.ErrIncorrectPassword) {
			message = err.Error()
		} else if e := password.Compare(id, got, want); e == nil {
			return "", nil
		}
	}
	return
}

func login(c *gin.Context) {
	var login struct {
		Username, Password string
		Rememberme         bool
	}
	if err := c.BindJSON(&login); err != nil {
		c.String(400, "")
		return
	}
	login.Username = strings.ToLower(login.Username)

	if password.IsMaxAttempts(info{login.Username, c.ClientIP()}) {
		c.JSON(200, gin.H{"status": 0, "message": fmt.Sprintf("Max retries exceeded (%d)", *maxRetry)})
		return
	}

	var user user
	var message string
	if err := accountClient.FindOne(mongodb.M{"username": login.Username}, nil, &user); err != nil {
		if err == mongodb.ErrNoDocuments {
			message = "Incorrect username"
		} else {
			svc.Print(err)
			c.String(500, "Critical Error! Please contact your system administrator.")
			return
		}
	} else {
		message, err = auth(info{login.Username, c.ClientIP()}, user.Password, login.Password)
		if err != nil {
			svc.Print(err)
			c.String(500, "Internal Server Error")
			return
		}

		if message == "" {
			session := sessions.Default(c)
			session.Clear()
			session.Set("id", user.ID.Hex())
			session.Set("username", user.Username)

			if login.Rememberme {
				session.Options(sessions.Options{HttpOnly: true, MaxAge: 856400 * 365})
			} else {
				session.Options(sessions.Options{HttpOnly: true})
			}

			if err := session.Save(); err != nil {
				svc.Print(err)
				c.String(500, "Internal Server Error")
				return
			}

			c.JSON(200, gin.H{"status": 1})
			return
		}
	}

	c.JSON(200, gin.H{"status": 0, "message": message})
}

func chgpwd(c *gin.Context) {
	session := sessions.Default(c)
	userID, username := session.Get("id"), session.Get("username")
	if userID == nil || username == nil {
		c.String(401, "")
		return
	}

	if password.IsMaxAttempts(info{username.(string), c.ClientIP()}) {
		c.JSON(200, gin.H{"status": 0, "message": fmt.Sprintf("Max retries exceeded (%d)", *maxRetry), "error": 1})
		return
	}

	var data struct{ Password, Password1, Password2 string }
	if err := c.BindJSON(&data); err != nil {
		c.String(400, "Bad Request")
		return
	}
	var err error
	if priv != nil {
		data.Password1, err = password.DecryptPKCS1v15(priv, data.Password1)
		if err != nil {
			c.String(400, "Bad Request")
			return
		}
		data.Password2, err = password.DecryptPKCS1v15(priv, data.Password2)
		if err != nil {
			c.String(400, "Bad Request")
			return
		}
	}

	id, _ := stockClient.ObjectID(userID.(string))
	var user user
	if err := accountClient.FindOne(mongodb.M{"_id": id}, nil, &user); err != nil {
		svc.Print(err)
		c.String(500, "")
		return
	}

	var message string
	var errorCode int
	if message, err = auth(info{username.(string), c.ClientIP()}, user.Password, data.Password); err != nil {
		svc.Print(err)
		c.String(500, "Internal Server Error")
		return
	} else {
		if priv != nil {
			data.Password, _ = password.DecryptPKCS1v15(priv, data.Password)
		}
		switch {
		case data.Password1 != data.Password2:
			message = "confirm password doesn't match new password"
			errorCode = 2
		case data.Password1 == data.Password:
			message = "new password cannot be the same as old password"
			errorCode = 2
		case data.Password1 == "":
			message = "new password cannot be blank"
		}
	}

	if message == "" {
		newPassword, err := password.HashPassword(data.Password1)
		if err != nil {
			svc.Print(err)
			c.String(500, "Internal Server Error")
			return
		}
		if _, err := accountClient.UpdateOne(
			mongodb.M{"_id": id},
			mongodb.M{"$set": mongodb.M{"password": newPassword}},
			nil,
		); err != nil {
			svc.Print(err)
			c.String(500, "")
			return
		}

		session.Clear()
		session.Options(sessions.Options{MaxAge: -1})
		if err := session.Save(); err != nil {
			svc.Print(err)
			c.String(500, "")
			return
		}

		c.JSON(200, gin.H{"status": 1})
		return
	} else {
		errorCode = 1
	}

	c.JSON(200, gin.H{"status": 0, "message": message, "error": errorCode})
}
