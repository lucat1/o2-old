package routes

import (
	"os"
	"path"

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
			"isownrepo":        shared.HasAccess(c, []string{"repo:settings"}, username, _repo.Name),
		})
		return
	} else {
		// Method POST
		owner := c.Param("user")
		reponame := c.Param("repo")

		newName := c.PostForm("general-rename")
		newOwner := c.PostForm("general-ownership")
		// Either it is rename or ownership change
		if (newName != "" && newName != reponame) || (newOwner != "" && newOwner != owner) {
			if newOwner == "" {
				newOwner = owner
			}
			if newName == "" {
				newName = reponame
			}

			shared.GetLogger().Info(
				"Renaming repository/Changing ownership",
				zap.String("oldOwner", owner),
				zap.String("newOwner", newOwner),
				zap.String("oldName", reponame),
				zap.String("newName", newName),
			)

			dbRepo := findRepoInDatabase(owner, reponame)
			err := shared.GetDatabase().Model(dbRepo).Update("name", newName).Update("owner", newOwner).Error
			if err != nil {
				shared.GetLogger().Warn("Could not rename repository", zap.Error(err))
				c.Status(500)
				c.Abort()
				return
			}

			os.MkdirAll(path.Join(cwd, "repos", newOwner), 0777)
			if err := os.Rename(getRepositoryPath(owner, reponame), getRepositoryPath(newOwner, newName)); err != nil {
				shared.GetLogger().Warn("Could not rename in the filesystem", zap.Error(err))
				c.Status(500)
				c.Abort()
				return
			}

			c.Redirect(302, "/"+newOwner+"/"+newName)
			return
		}

		// Delete the repo
		if c.PostForm("general-delete") == "true" {
			shared.GetLogger().Info(
				"Deleting repository",
				zap.String("owner", owner),
				zap.String("repo", reponame),
			)

			dbRepo := findRepoInDatabase(owner, reponame)
			if err := shared.GetDatabase().Delete(dbRepo).Error; err != nil {
				shared.GetLogger().Warn("Could not reomve repository from database", zap.Error(err))
				c.Status(500)
				c.Abort()
				return
			}

			if err := os.RemoveAll(getRepositoryPath(owner, reponame)); err != nil {
				shared.GetLogger().Warn("Could not reomve repository from filesystem", zap.Error(err))
				c.Status(500)
				c.Abort()
				return
			}

			c.Redirect(301, "/"+owner)
			return
		}

		NotFound(c)
		c.Abort()
		return
	}
}
