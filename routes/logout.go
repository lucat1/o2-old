package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Logout removes the session token from the request
// /logout
func Logout(c *gin.Context) {
	// Since we're listening for /:user we must
	// Check if the parameter is logout
	if c.Param("user") != "logout" {
		c.Next() // Skip
		return
	}

	if c.Keys["user"] == nil {
		c.Redirect(301, "/login")
		c.Abort()
		return
	}

	// Remove the cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   "invalid",
		Expires: time.Now().Add(time.Hour * 24),
		Path:    "/",
	})
	c.Redirect(301, "/")
	c.Abort()
}
