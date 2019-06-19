package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

func ExistsRepo(withDatabase bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("user")
		reponame := c.Param("repo")

		if c.Keys == nil {
			c.Keys = make(map[string]interface{})
		}

		if withDatabase {
			_, repo := findRepo(c, username, reponame)
			if repo == nil {
				NotFound(c)
				return
			}

			c.Keys["repo"] = repo
			c.Next()
		} else {
			// Used for git push / pull
			// Faster without database, for lots of consequent requests
			dir := getRepositoryPath(username, reponame)
			_, err := os.Stat(dir)
			if err != nil {
				fmt.Println(err)
				c.Status(404)
				return
			}

			c.Keys["dir"] = dir
			c.Keys["file"] = c.Param("path")
			c.Next()
		}
	}
}