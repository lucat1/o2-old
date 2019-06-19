package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lucat1/git/shared"
	"go.uber.org/zap"
)

// Register serves the rendered page for login
// /register
func Register(c *gin.Context) {
	// Register tough has to first check out we're not
	// in a user path like /luca but we are in fact in /register
	if c.Request.URL.Path != "/register" {
		c.Next() // Skip
		return
	}

	if c.Request.Method == "GET" {
		c.HTML(200, "register.tmpl", gin.H{
			"user": c.Keys["user"],
		})
		c.Abort()
	} else {
		username := c.PostForm("username")
		password := c.PostForm("password")
		shared.GetLogger().Info(
			"New registration",
			zap.String("username", username),
			zap.String("password", password),
		)
		passwd, err := shared.HashPassword(password)
		if err != nil {
			shared.GetLogger().Error(
				"Could not hash password",
				zap.String("password", password),
				zap.Error(err),
			)
			c.HTML(500, "register.tmpl", gin.H{
				"error":   true,
				"message": "Could not hash your password",
			})
			return
		}

		// Create the user and add it to the database
		user := shared.User{
			Username: username,
			Password: passwd,
		}
		err = shared.GetDatabase().Save(&user).Error
		if err != nil {
			c.HTML(500, "register.tmpl", gin.H{
				"error":   true,
				"message": "Could not create user, does it already exist?",
			})
			return
		}
		authenticate(c, user)
		c.Redirect(301, "/"+username)
	}
}
