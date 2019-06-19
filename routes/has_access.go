package routes

import "github.com/gin-gonic/gin"

func hasAccess(c *gin.Context, scopes []string) bool {
	// Ignore unauthenticated requests
	if c.Keys["user"] == nil {
		return false
	}

	// TODO: Check if user has access to scopes

	return true
}

// HasAccess checks if the authenticated user
// has access to a certain repository feature
func HasAccess(scopes []string) func(*gin.Context) {
	return func(c *gin.Context) {
		if hasAccess(c, scopes) {
			c.Next()
		} else {
			NotFound(c)
		}
	}
}
