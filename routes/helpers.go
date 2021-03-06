package routes

import (
	"os"
	"path"

	"code.gitea.io/git"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/lucat1/o2/shared"
	"go.uber.org/zap"
)

var cwd string

func getRepositoryPath(user, repo string) string {
	var wd string
	if cwd == "" {
		_wd, err := os.Getwd()
		if err != nil {
			shared.GetLogger().Fatal("Could not get cwd via syscall", zap.Error(err))
		}
		cwd = _wd
		wd = _wd
	} else {
		wd = cwd
	}

	return path.Join(wd, "repos", user, repo+".git")
}

func getRepository(c *gin.Context, user, repo string) *git.Repository {
	path := getRepositoryPath(user, repo)

	r, err := git.OpenRepository(path)
	if err != nil {
		NotFound(c)
	}

	return r
}

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
