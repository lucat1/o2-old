package routes

import (
	"html/template"
	"io/ioutil"
	"net/http"

	"code.gitea.io/git"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/jinzhu/gorm"
	"github.com/lucat1/git/shared"
	"go.uber.org/zap"
)

func findRepoInDatabase(username string, reponame string) *shared.Repository {
	var repo shared.Repository
	err := shared.GetDatabase().Find(&repo, &shared.Repository{Name: reponame, Owner: username}).Error
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			shared.GetLogger().Error(
				"Unknown error in db while finding repository",
				zap.String("username", username),
				zap.String("reponame", reponame),
				zap.Error(err),
			)
		}
		return nil
	}
	return &repo
}

func findRepo(c *gin.Context, username string, reponame string) (*shared.Repository, *git.Repository) {
	_repo := findRepoInDatabase(username, reponame)
	if _repo == nil {
		NotFound(c)
		return nil, nil
	}
	repo := getRepository(c, username, reponame)
	if repo == nil {
		NotFound(c)
		return _repo, nil
	}

	return _repo, repo
}

// Repo renders the repository view
// /:user/:repo
func Repo(c *gin.Context) {
	username := c.Param("user")
	reponame := c.Param("repo")
	_repo, repo := findRepo(c, username, reponame)
	if repo == nil {
		return
	}

	commit := getCommit(c, repo, _repo.MainBranch)

	if commit == nil {
		// Uninitialized repo, should simply display clone infos
		c.HTML(http.StatusOK, "repo.tmpl", gin.H{
			"username":     username,
			"user":         c.Keys["user"],
			"repo":         reponame,
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
			AbsolutePrefix: "http://o2.local/" + username + "/" + reponame + "/blob/master",
		}))
	}

	c.HTML(http.StatusOK, "repo.tmpl", gin.H{
		"username":     username,
		"user":         c.Keys["user"],
		"repo":         reponame,
		"selectedrepo": true,
		"markdown":     template.HTML(md),
	})
}
