package routes

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/lucat1/git/shared"

	"code.gitea.io/git"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func getCommit(c *gin.Context, repo *git.Repository, branch string) *git.Commit {
	var (
		commit *git.Commit
		err    error
	)

	// Find by commit SHA not by branch
	if len(branch) == 40 || len(branch) == 20 {
		commit, err = repo.GetCommit(branch)
	} else {
		commit, err = repo.GetBranchCommit(branch)
	}
	if err != nil {
		return nil
	}

	return commit
}

// Entry represents a file in the tree
type Entry struct {
	Mode    string
	Name    string
	IsDir   bool
	Size    string
	Summary string
}

func stringifyMode(mode git.EntryMode) string {
	switch mode {
	case git.EntryModeBlob:
		return "-rw-r--r--"
	case git.EntryModeTree:
		return "d---------"
	case git.EntryModeExec:
		return "-rwxr-xr-x"
	case git.EntryModeSymlink:
		return "lrwxr-xr-x"
	case git.EntryModeCommit:
		return "----------"
	}
	return ""
}

func buildEntry(commitAndEntry []interface{}) *Entry {
	entry := commitAndEntry[0].(*git.TreeEntry)
	commit := commitAndEntry[1].(*git.Commit)

	return &Entry{
		Mode:    stringifyMode(entry.Mode()),
		Name:    entry.Name(),
		IsDir:   entry.IsDir(),
		Size:    humanize.Bytes(uint64(entry.Size())),
		Summary: commit.Summary(),
	}
}

// Tree renders the repository tree
// /:user/:repo/tree
func Tree(c *gin.Context) {
	username := c.Param("user")
	ref := c.Param("ref")
	path := c.Param("path")
	if path != "" {
		path = strings.Replace(path, "/", "", 1)
	}

	_Irepo, Irepo := c.Keys["_repo"], c.Keys["repo"]
	if Irepo == nil || _Irepo == nil {
		NotFound(c)
		return
	}

	dbRepo := _Irepo.(*shared.Repository)
	repo := Irepo.(*git.Repository)

	if ref == "" {
		ref = dbRepo.MainBranch
	}

	commit := getCommit(c, repo, ref)
	if commit == nil {
		NotFound(c)
		return
	}

	var (
		entry  interface{}
		err    error
		isTree bool
	)

	if path == "/" {
		entry = &commit.Tree
		isTree = true
	} else {
		entry, err = commit.GetTreeEntryByPath(path)
	}

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	// Paths for the separator
	var parents [][]string
	prev := ""
	_parents := strings.Split(path, "/")
	for _, path := range _parents {
		p := path
		if path == "" {
			p = dbRepo.Name
		} else {
			prev += "/" + path
		}
		parents = append(parents, []string{p, prev})
	}

	_path := path
	if path != "" {
		_path = "/" + path
	}

	if isTree || entry.(*git.TreeEntry).IsDir() {
		sub, _ := commit.SubTree(path)
		entries, _ := sub.ListEntries()
		entries.Sort()

		_files, err := entries.GetCommitsInfo(commit, path, nil)
		if err != nil {
			shared.GetLogger().Error("Could not load commits infos", zap.Error(err))
			NotFound(c)
			return
		}

		var files []*Entry
		for _, file := range _files {
			files = append(files, buildEntry(file))
		}

		c.HTML(http.StatusOK, "tree.tmpl", gin.H{
			"username":     username,
			"user":         c.Keys["user"],
			"repo":         dbRepo.Name,
			"isownrepo":    isOwnRepo(c, dbRepo.Owner),
			"ref":          ref,
			"path":         _path,
			"parents":      parents,
			"parentsl":     len(parents) - 1,
			"selectedtree": true,
			"directory":    true,
			"notRoot":      !(path == "/"),
			"files":        files,
		})
	} else {
		entry := entry.(*git.TreeEntry)
		reader, err := entry.Blob().Data()
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		b, err := ioutil.ReadAll(reader)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		var contents interface{} // either string or tempalte.HTML

		lexer := lexers.Match(entry.Name())
		if lexer == nil {
			shared.GetLogger().Info(
				"Could not get language for file",
				zap.String("filename", entry.Name()),
			)
			lexer = lexers.Fallback
		}
		formatter := html.New(
			html.WithLineNumbers(),
			html.Standalone(),
			html.WithClasses(),
			html.TabWidth(2),
		)
		style := styles.Get("xcode")
		iterator, err := lexer.Tokenise(nil, string(b))
		var buf bytes.Buffer
		err = formatter.Format(&buf, style, iterator)
		if err != nil {
			contents = string(b)
		} else {
			contents = template.HTML(buf.String())
		}

		//data := buildEntry(entry)
		c.HTML(http.StatusOK, "tree.tmpl", &gin.H{
			"username":     username,
			"user":         c.Keys["user"],
			"repo":         dbRepo.Name,
			"isownrepo":    isOwnRepo(c, dbRepo.Owner),
			"ref":          ref,
			"path":         _path,
			"parents":      parents,
			"parentsl":     len(parents) - 1,
			"selectedtree": true,
			"directory":    false,
			"contents":     contents,
			"mode":         stringifyMode(entry.Mode()),
			"size":         humanize.Bytes(uint64(entry.Size())),
		})
	}
}
