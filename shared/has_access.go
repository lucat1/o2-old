package shared

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

// HasAccess cheks if the currently authenticated user has
// access to the provided scopes on a source(identified by owner and repo)
func HasAccess(c *gin.Context, scopes []string, owner string, repo string) bool {
	// Ignore unauthenticated requests
	if c.Keys["user"] == nil {
		return false
	}

	user := c.Keys["user"].(*User)

	GetLogger().Info(
		"Authorization request",
		zap.String("asking", user.Username),
		zap.Strings("resource", []string{owner, repo}),
		zap.Strings("scopes", scopes),
	)

	var dbRepo Repository
	err := GetDatabase().Find(&dbRepo, "owner = ? AND name = ?", owner, repo).Error
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			GetLogger().Warn("Unkown error while checking user's access", zap.Error(err))
		}

		return false
	}

	restricted := strings.Split(dbRepo.Restricted, ",")

	for _, scope := range scopes {
		for _, _scope := range restricted {
			if scope == _scope {
				// If the required scope is restricted,
				// then check that the permission is granted to the user
				ok := false
				permissions := ParsePermissions(&dbRepo)
				for _, granted := range *permissions {
					if granted.For == user.Username && granted.Key == scope {
						ok = true
						break
					}
				}

				if !ok {
					return false
				}
			}
		}
	}

	return true
}
