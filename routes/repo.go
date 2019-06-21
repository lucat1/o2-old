package routes

import (
	"html/template"
	"io/ioutil"
	"net/http"

	"code.gitea.io/git"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/lucat1/o2/shared"
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
	var md []byte

	if commit == nil {
		// No commits yet, prompt the user to clone the repo and push
		md = []byte("Please clone the repo and start committing to it")
	} else {
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
				AbsolutePrefix: "/" + username + "/" + _repo.Name,
			}))
		}
	}

	c.HTML(http.StatusOK, "repo.tmpl", gin.H{
		"username":     username,
		"user":         c.Keys["user"],
		"repo":         _repo.Name,
		"selectedrepo": true,
		"isownrepo":    isOwnRepo(c, _repo.Owner),
		"markdown":     template.HTML(md),
	})
}
