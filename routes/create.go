package routes

import (
	"code.gitea.io/git"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/lucat1/git/shared"
)

// Create is the /create route, to create repos
// /create
func Create(c *gin.Context) {
	// Create tough has to first check out we're not
	// in a user path like /luca but we are in fact in /create
	if c.Request.URL.Path != "/create" {
		c.Next() // Skip
		return
	}

	// Check if the user is authenticated first
	if c.Keys == nil || c.Keys["user"] == nil {
		c.Redirect(301, "/login")
		return
	}

	if c.Request.Method == "GET" {
		// Check if the user is authenticated first
		c.HTML(200, "create.tmpl", gin.H{
			"user": c.Keys["user"],
		})
		c.Abort()
	} else {
		// POST
		user := c.Keys["user"].(*shared.User)
		name := c.PostForm("name")

		var oldRepo shared.Repository
		err := shared.GetDatabase().Find(&oldRepo, &shared.Repository{Owner: user.Username, Name: name}).Error
		if !gorm.IsRecordNotFoundError(err) {
			failWithError(c, "You already own a repository with this name")
			return
		}
		// All good to go
		repo := &shared.Repository{
			Owner:      user.Username,
			Name:       name,
			MainBranch: "master", // TODO: Form parameter
		}
		err = shared.GetDatabase().Save(repo).Error
		if err != nil {
			failWithError(c, err.Error())
			return
		}

		err = git.InitRepository(getRepositoryPath(user.Username, name), true)
		if err != nil {
			failWithError(c, err.Error())
			return
		}

		c.Redirect(301, "/"+user.Username+"/"+name)
	}
}

func failWithError(c *gin.Context, err string) {
	c.HTML(200, "create.tmpl", gin.H{
		"user":  c.Keys["user"],
		"error": err,
	})
	c.Abort()
}
