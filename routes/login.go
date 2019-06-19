package routes

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/lucat1/git/shared"
	"go.uber.org/zap"
)

func authenticate(c *gin.Context, user shared.User) {
	expirationTime := time.Now().Add(24 * time.Hour)
	// Create the JWT claims, which includes the username and expiry time
	claims := &shared.Claims{
		UUID: user.ID.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(shared.JWT)
	if err != nil {
		shared.GetLogger().Error("Could not convert JWT into string", zap.Error(err))
		NotFound(c)
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
		Path:    "/",
	})
}

// Login serves the rendered page for login
// /login
func Login(c *gin.Context) {
	// Login tough has to first check out we're not
	// in a user path like /luca but we are in fact in /login
	if c.Param("user") != "login" {
		c.Next() // Skip
		return
	}

	if c.Request.Method == "GET" {
		if c.Keys["user"] != nil {
			c.Redirect(301, "/"+c.Keys["user"].(*shared.User).Username)
			c.Abort()
			return
		}

		c.HTML(200, "login.tmpl", gin.H{
			"user": c.Keys["user"],
		})
		c.Abort()
	} else {
		username := c.PostForm("username")
		password := c.PostForm("password")
		shared.GetLogger().Info(
			"New login",
			zap.String("username", username),
			zap.String("password", password),
		)

		user := FindUser(username)
		if user == nil {
			c.HTML(500, "login.tmpl", gin.H{
				"error":   true,
				"message": "Could not find the user",
			})
			return
		}

		ok := shared.CheckPassword(user.Password, password)
		if ok {
			authenticate(c, *user)
			c.Redirect(301, "/"+user.Username)
		} else {
			c.HTML(500, "login.tmpl", gin.H{
				"error":   true,
				"message": "Invalid password",
			})
			return
		}
	}
}
