package routes

import (
	"fmt"
	"net/http"

	"code.gitea.io/git"
	"github.com/gin-gonic/gin"
	"github.com/lucat1/git/shared"
)

type renderCommit struct {
	ID          string
	Name        string
	Description string
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
		commits = append(commits, &renderCommit{
			ID:          commit.ID.String(),
			Description: commit.Summary(),
			Name:        commit.CommitMessage,
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
