package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucat1/o2/shared"
	"go.uber.org/zap"
)

// HasAccess checks if the authenticated user
// has access to a certain repository feature
//
// Must have path:
// /:user/:repo/WHATEVER
func HasAccess(scopes []string) func(*gin.Context) {
	return func(c *gin.Context) {
		if shared.HasAccess(c, scopes, c.Param("user"), c.Param("repo")) {
			c.Next()
		} else {
			NotFound(c)
		}
	}
}

// RawHasAccess handles protected routes in a raw way, not displaying
// a user-friendly 404 but returning a Unauthorized status, also prompting
// for basic auth instead of using JWTs
//
// This is used only in git push/pull routes
//
// Must have path:
// /:user/:repo/WHATEVER
func RawHasAccess(scopes []string) func(*gin.Context) {
	return func(c *gin.Context) {
		authHead := c.GetHeader("Authorization")
		if len(authHead) == 0 {
			c.Header("WWW-Authenticate", "Basic realm=\".\"")
			c.Status(http.StatusUnauthorized)
			return
		}

		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.Status(http.StatusBadRequest)
			return
		}

		shared.GetLogger().Info(
			"New login",
			zap.String("username", username),
			zap.String("password", password),
		)

		user := FindUser(username)
		if user == nil {
			// User with the provided username doesnt exist
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}
		// If the user exists lets check for the password
		ok = shared.CheckPassword(user.Password, password)
		if !ok {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		shared.GetLogger().Info("Raw user authenticated", zap.String("user", user.Username))

		c.Keys["user"] = user
		if ok := shared.HasAccess(c, scopes, c.Param("user"), c.Param("repo")); !ok {
			c.Status(http.StatusForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}
