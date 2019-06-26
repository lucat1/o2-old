package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lucat1/o2/shared"
	"go.uber.org/zap"
)

// Settings renders the settings page for a repository
// /:user/:repo/settings
func Settings(c *gin.Context) {
	if c.Request.Method == "GET" {
		username := c.Param("user")
		_Irepo, Irepo := c.Keys["_repo"], c.Keys["repo"]
		if Irepo == nil || _Irepo == nil {
			NotFound(c)
			return
		}

		_repo := _Irepo.(*shared.Repository)

		c.HTML(200, "settings.tmpl", gin.H{
			"username":         username,
			"user":             c.Keys["user"],
			"repo":             _repo.Name,
			"mainbranch":       _repo.MainBranch,
			"selectedsettings": true,
			"isownrepo":        isOwnRepo(c, _repo.Owner),
		})
		return
	} else {
		owner := c.Param("user")
		reponame := c.Param("repo")
		// Method POST
		if newName := c.PostForm("general-rename"); newName != "" && newName != reponame {
			shared.GetLogger().Info(
				"Renaming repository",
				zap.String("owner", owner),
				zap.String("from", reponame),
				zap.String("to", newName),
			)

			dbRepo := findRepoInDatabase(owner, reponame)
			if err := shared.GetDatabase().Model(dbRepo).Update("name", newName).Error; err != nil {
				shared.GetLogger().Warn(
					"Could not rename repository",
					zap.String("owner", owner),
					zap.String("from", reponame),
					zap.String("to", newName),
					zap.Error(err),
				)
				c.Status(500)
				c.Abort()
				return
			}
			if err := os.Rename(getRepositoryPath(owner, reponame), getRepositoryPath(owner, newName)); err != nil {
				shared.GetLogger().Warn(
					"Could not rename in the filesystem",
					zap.String("owner", owner),
					zap.String("from", reponame),
					zap.String("to", newName),
					zap.Error(err),
				)
				c.Status(500)
				c.Abort()
				return
			}

			c.Redirect(302, "/"+owner+"/"+newName)
			return
		}

	}
}
