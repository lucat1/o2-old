package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lucat1/o2/shared"
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
	}
}
