package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/lucat1/git/shared"
	"go.uber.org/zap"
)

func findUser(username string) *shared.User {
	var user shared.User
	err := shared.GetDatabase().Find(&user, &shared.User{Username: username}).Error
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			shared.GetLogger().Error(
				"Unknown error in db while finding user",
				zap.String("username", username),
				zap.Error(err),
			)
		}
		return nil
	}

	return &user
}

// User renders the user's profile
// /:user
func User(c *gin.Context) {
	username := c.Param("user")
	user := findUser(username)
	if user == nil {
		NotFound(c)
		return
	}

	var repos []*shared.Repository
	if err := shared.GetDatabase().Find(&repos, &shared.Repository{ Owner: user.Username }).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		shared.GetLogger().Error(
			"Unkown error while listing user's repositories",
			zap.String("username", user.Username),
			zap.Error(err),
		)
	}

	c.HTML(200, "user.tmpl", gin.H{
		"username":    username,
		"email":       user.Email,
		"firstname":   user.Firstname,
		"lastname":    user.Lastname,
		"description": user.Description,
		"repos":       repos,
		"user":        c.Keys["user"],
	})
}
