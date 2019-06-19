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

func buildEntry(entry *git.TreeEntry) []interface{} {
	var mode string
	switch entry.Mode() {
	case git.EntryModeBlob:
		mode = "-rw-r--r--"
	case git.EntryModeTree:
		mode = "d---------"
	case git.EntryModeExec:
		mode = "-rwxr-xr-x"
	case git.EntryModeSymlink:
		mode = "lrwxr-xr-x"
	case git.EntryModeCommit:
		mode = "----------"
	}

	return []interface{}{
		mode,
		entry.Name(),
		entry.IsDir(),
		humanize.Bytes(uint64(entry.Size())),
	}
}

// Tree renders the repository tree
// /:user/:repo/tree
func Tree(c *gin.Context) {
	username := c.Param("user")
	ref := c.Param("ref")
	path := c.Param("path")
	if path == "" {
		path = "/"
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

	if isTree || entry.(*git.TreeEntry).IsDir() {
		sub, _ := commit.SubTree(path)
		entries, _ := sub.ListEntries()
		var files [][]interface{}
		for _, entry := range entries {
			files = append(files, buildEntry(entry))
		}
		c.HTML(http.StatusOK, "tree.tmpl", gin.H{
			"username":     username,
			"user":         c.Keys["user"],
			"repo":         dbRepo.Name,
			"ref":          ref,
			"path":         path,
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

		data := buildEntry(entry)
		c.HTML(http.StatusOK, "tree.tmpl", &gin.H{
			"username":     username,
			"user":         c.Keys["user"],
			"repo":         dbRepo.Name,
			"ref":          ref,
			"path":         path,
			"parents":      parents,
			"parentsl":     len(parents) - 1,
			"selectedtree": true,
			"directory":    false,
			"contents":     contents,
			"mode":         data[0].(string),
			"size":         data[3],
		})
	}
}
