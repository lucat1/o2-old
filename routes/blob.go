package routes

import (
	"io"

	"code.gitea.io/git"
	"github.com/gin-gonic/gin"
)

// Blob reponds with the raw blob of a file
// /:user/:repo/blob
func Blob(c *gin.Context) {
	ref := c.Param("ref")
	path := c.Param("path")

	_, Irepo := c.Keys["_repo"], c.Keys["repo"]
	if Irepo == nil {
		NotFound(c)
		return
	}

	repo := Irepo.(*git.Repository)

	commit := getCommit(c, repo, ref)

	blob, err := commit.GetBlobByPath(path)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	reader, err := blob.Data()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.Header("Content-Type", "text/plain")
	io.Copy(c.Writer, reader)
}
