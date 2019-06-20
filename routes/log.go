package routes

import (
	"fmt"
	"net/http"
	"strings"

	"code.gitea.io/git"
	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"github.com/lucat1/git/shared"
)

type renderCommit struct {
	ID          string
	ShortID     string
	Author      *git.Signature
	Name        string
	Description string
	Time        string
}

// Log renders the repository commits
// /:user/:repo/log
func Log(c *gin.Context) {
	username := c.Param("user")
	_Irepo, Irepo := c.Keys["_repo"], c.Keys["repo"]
	if Irepo == nil || _Irepo == nil {
		NotFound(c)
		return
	}

	dbRepo := _Irepo.(*shared.Repository)
	repo := Irepo.(*git.Repository)

	commit := getCommit(c, repo, dbRepo.MainBranch)
	if commit == nil {
		NotFound(c)
		return
	}
	_commits, err := commit.CommitsBeforeLimit(20)
	if err != nil {
		c.AbortWithError(500, fmt.Errorf("Could not load commits %e", err))
		return
	}

	var commits []*renderCommit
	for e := _commits.Front(); e != nil; e = e.Next() {
		commit := e.Value.(*git.Commit)
		id := commit.ID.String()
		commits = append(commits, &renderCommit{
			ID:          id,
			ShortID:     id[:8],
			Author:      commit.Author,
			Name:        commit.Summary(),
			Description: strings.Replace(commit.Message(), commit.Summary()+"\n", "", 1),
			Time:        humanize.Time(commit.Author.When),
		})
	}

	c.HTML(http.StatusOK, "log.tmpl", gin.H{
		"username":    username,
		"repo":        dbRepo.Name,
		"isownrepo":   isOwnRepo(c, dbRepo.Owner),
		"selectedlog": true,
		"commits":     commits,
		"user":        c.Keys["user"],
	})
}
