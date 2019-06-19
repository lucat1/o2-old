package routes

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/lucat1/git/shared"
	uuid "github.com/satori/go.uuid"
)

// AuthMiddleware sets the context.key["user"] value
// For all other requests to use
func AuthMiddleware(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}

	if err == nil {
		claims := &shared.Claims{}
		tkn, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
			return shared.JWT, nil
		})
		if err == nil && tkn.Valid {
			var user shared.User
			shared.GetDatabase().Find(&user, &shared.User{ID: uuid.FromStringOrNil(claims.UUID)})

			// If we're at 1 minute(or lower) from the expiry, refresh the token
			if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < time.Hour {
				shared.GetLogger().Info("Refreshing token")
				authenticate(c, user)
			}
			c.Keys["user"] = &user
		}
	}

	c.Next()
}

// WithAuth is a middleware that checks for authentication
func WithAuth(h gin.HandlerFunc, scopes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Keys["user"] != nil {

		}
	}
}
