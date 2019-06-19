package routes

import (
	"html/template"
	"io/ioutil"
	"net/http"

	"code.gitea.io/git"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/lucat1/git/shared"
)

// Repo renders the repository view
// /:user/:repo
func Repo(c *gin.Context) {
	username := c.Param("user")
	_Irepo, Irepo := c.Keys["_repo"], c.Keys["repo"]
	if Irepo == nil || _Irepo == nil {
		NotFound(c)
		return
	}

	_repo := _Irepo.(*shared.Repository)
	repo := Irepo.(*git.Repository)

	commit := getCommit(c, repo, _repo.MainBranch)

	if commit == nil {
		// Uninitialized repo, should simply display clone infos
		c.HTML(http.StatusOK, "repo.tmpl", gin.H{
			"username":     username,
			"user":         c.Keys["user"],
			"repo":         _repo.Name,
			"selectedrepo": true,
			"markdown":     "Please clone the repo and start committing to it",
		})
		return
	}

	var md []byte
	blob, err := commit.GetBlobByPath("README.md")
	if err != nil {
		md = []byte{}
	} else {
		reader, err := blob.Data()
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		bytes, err := ioutil.ReadAll(reader)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		md = markdown.ToHTML(bytes, nil, html.NewRenderer(html.RendererOptions{
			AbsolutePrefix: "http://o2.local/" + username + "/" + _repo.Name + "/blob/master",
		}))
	}

	c.HTML(http.StatusOK, "repo.tmpl", gin.H{
		"username":     username,
		"user":         c.Keys["user"],
		"repo":         _repo.Name,
		"selectedrepo": true,
		"markdown":     template.HTML(md),
	})
}
