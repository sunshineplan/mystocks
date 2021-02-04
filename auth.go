package main

import (
	"database/sql"
	"log"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	ID       int
	Username string
	Password string
}

func authRequired(c *gin.Context) {
	if sessions.Default(c).Get("id") == nil {
		c.Abort()
		c.Redirect(302, "/")
	}
}

func getUser(c *gin.Context) (id int, username string, err error) {
	var ok bool
	id, ok = sessions.Default(c).Get("id").(int)
	if !ok {
		return
	}
	if universal {
		err = db.QueryRow("SELECT id FROM user WHERE uid = ?", id).Scan(&id)
		if err != nil {
			return
		}
		username = sessions.Default(c).Get("username").(string)
		return
	}
	err = db.QueryRow("SELECT username FROM user WHERE id = ?", id).Scan(&username)
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

	var user user
	statusCode := 200
	var message string
	if err := db.QueryRow(
		"SELECT id, username, password FROM user WHERE username = ?",
		login.Username,
	).Scan(&user.ID, &user.Username, &user.Password); err != nil {
		if strings.Contains(err.Error(), "doesn't exist") {
			restore("")
			statusCode = 503
			message = "Detected first time running. Initialized the database."
		} else if err == sql.ErrNoRows {
			statusCode = 403
			message = "Incorrect username"
		} else {
			log.Print(err)
			statusCode = 500
			message = "Critical Error! Please contact your system administrator."
		}
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
			if (err == bcrypt.ErrHashTooShort && user.Password != login.Password) ||
				err == bcrypt.ErrMismatchedHashAndPassword {
				statusCode = 403
				message = "Incorrect password"
			} else if user.Password != login.Password {
				log.Print(err)
				statusCode = 500
				message = "Critical Error! Please contact your system administrator."
			}
		}
		if message == "" {
			session := sessions.Default(c)
			session.Clear()
			session.Set("id", user.ID)

			if login.Rememberme {
				session.Options(sessions.Options{Path: "/", HttpOnly: true, MaxAge: 856400 * 365})
			} else {
				session.Options(sessions.Options{Path: "/", HttpOnly: true})
			}

			if err := session.Save(); err != nil {
				log.Print(err)
				statusCode = 500
				message = "Failed to save session."
			}
		}
	}
	c.String(statusCode, message)
}

func chgpwd(c *gin.Context) {
	var data struct{ Password, Password1, Password2 string }
	if err := c.BindJSON(&data); err != nil {
		c.String(400, "")
		return
	}

	session := sessions.Default(c)
	userID := session.Get("id")

	var oldPassword string
	if err := db.QueryRow("SELECT password FROM user WHERE id = ?", userID).Scan(&oldPassword); err != nil {
		log.Print(err)
		c.String(500, "")
		return
	}

	var message string
	var errorCode int
	err := bcrypt.CompareHashAndPassword([]byte(oldPassword), []byte(data.Password))
	switch {
	case err != nil:
		if (err == bcrypt.ErrHashTooShort && data.Password != oldPassword) ||
			err == bcrypt.ErrMismatchedHashAndPassword {
			message = "Incorrect password."
			errorCode = 1
		} else if data.Password != oldPassword {
			log.Print(err)
			c.String(500, "")
			return
		}
	case data.Password1 != data.Password2:
		message = "Confirm password doesn't match new password."
		errorCode = 2
	case data.Password1 == data.Password:
		message = "New password cannot be the same as your current password."
		errorCode = 2
	case data.Password1 == "":
		message = "New password cannot be blank."
	}

	if message == "" {
		newPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password1), bcrypt.MinCost)
		if err != nil {
			log.Print(err)
			c.String(500, "")
			return
		}
		if _, err := db.Exec("UPDATE user SET password = ? WHERE id = ?", string(newPassword), userID); err != nil {
			log.Print(err)
			c.String(500, "")
			return
		}
		session.Clear()
		if err := session.Save(); err != nil {
			log.Print(err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		return
	}
	c.JSON(200, gin.H{"status": 0, "message": message, "error": errorCode})
}
