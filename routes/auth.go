package routes

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/lucat1/o2/shared"
	uuid "github.com/satori/go.uuid"
)

// AuthMiddleware sets the context.Key["user"] value
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

			// If we're at 1 hour(or lower) from the expiry date, refresh the token
			if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < time.Hour {
				shared.GetLogger().Info("Refreshing token")
				authenticate(c, user)
			}
			c.Keys["user"] = &user
		}
	}

	c.Next()
}
