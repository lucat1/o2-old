package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucat1/git/shared"
)

// Logout removes the session token from the request
// /logout
func Logout(c *gin.Context) {
	// Since we're listening for /:user we must
	// Check if the parameter is logout
	if c.Request.URL.Path != "/logout" {
		c.Next() // Skip
		return
	}

	if c.Keys["user"] == nil {
		return
	}

	shared.GetLogger().Info("Updating cookie")

	// Remove the cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:  "token",
		Value: "",
		Path:  "/",
	})
	c.Redirect(301, "/")
	c.Abort()
}
