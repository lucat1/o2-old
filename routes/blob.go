package routes

import (
	"io"

	"github.com/gin-gonic/gin"
)

// Blob reponds with the raw blob of a file
// /:user/:repo/blob
func Blob(c *gin.Context) {
	user := c.Param("user")
	_repo := c.Param("repo")
	ref := c.Param("ref")
	path := c.Param("path")

	repo := getRepository(c, user, _repo)
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
