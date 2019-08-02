package routes

import (
	"html/template"
	"path/filepath"
	"strconv"
	"strings"

	"code.gitea.io/git"
	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"github.com/lucat1/o2/shared"
	"go.uber.org/zap"
)

type LineType uint8

const (
	LineAddition  LineType = 0
	LineDeletion  LineType = 1
	LineUnchanged LineType = 2
)

type FileDiff struct {
	FromFile string
	ToFile   string
	Parts    []*Part
	Change   string
}

type Part struct {
	Header string
	Start  int

	LatestAdded   int
	LatestRemoved int
	Lines         []*Line
}

type Line struct {
	raw string

	OldNumber int
	NewNumber int
	Type      LineType
}

func (l *Line) Render() template.HTML {
	return template.HTML(strings.Replace(strings.Replace(l.raw, "+", "", 1), "-", "", 1))
}

func (l *Line) Class() string {
	switch l.Type {
	case LineAddition:
		return "added"
	case LineDeletion:
		return "deleted"
	}

	// Never hit
	return ""
}

func getLatestPart(diffs []*FileDiff) *Part {
	diff := diffs[len(diffs)-1]
	if len(diff.Parts) == 0 {
		return nil
	}

	return diff.Parts[len(diff.Parts)-1]
}

// Diff renders a commit diff
// /:user/:repo/diff/:sha
func Diff(c *gin.Context) {
	username := c.Param("user")
	sha := c.Param("sha")
	_Irepo, Irepo := c.Keys["_repo"], c.Keys["repo"]
	if Irepo == nil || _Irepo == nil {
		NotFound(c)
		return
	}

	_repo := _Irepo.(*shared.Repository)
	repo := Irepo.(*git.Repository)
	repoPath := getRepositoryPath(username, c.Param("repo"))

	// Find this commit and the previous one
	commit := getCommit(c, repo, sha)
	firstID := commit.ID.String()
	commits, err := commit.CommitsBeforeLimit(2)
	if err != nil || commits.Front() == nil {
		shared.GetLogger().Warn("Error while getting previos commit in diff", zap.String("id", firstID), zap.Error(err))
		NotFound(c)
		return
	}
	second := commits.Front().Next()
	if second == nil {
		shared.GetLogger().Warn("Diffing first commit is not yet supported", zap.String("id", firstID))
		NotFound(c)
		return
	}
	secondID := second.Value.(*git.Commit).ID.String()

	statuses, err := git.GetCommitFileStatus(repoPath, firstID)
	if err != nil {
		shared.GetLogger().Error("Could not diff commit", zap.String("id", firstID), zap.Error(err))
		NotFound(c)
		return
	}

	args := append([]string{"diff", secondID, firstID, "--"}, statuses.Modified...)
	cmd := git.NewCommand(args...)
	out, err := cmd.RunInDir(repoPath)
	if err != nil {
		shared.GetLogger().Error("Could not run git diff", zap.Strings("args", args), zap.Error(err))
		NotFound(c)
		return
	}

	insertions, deletions := 0, 0
	var diffs []*FileDiff
	var latestDiffBeginning = 0
	lines := strings.Split(out, "\n")
	for i, line := range lines {
		// Beginning of a file diff
		if strings.HasPrefix(line, "diff --git ") {
			files := strings.Split(strings.Replace(line, "diff --git ", "", 1), " ")
			fromFile, toFile := files[0], files[1]
			diffs = append(diffs, &FileDiff{
				FromFile: filepath.Base(fromFile),
				ToFile:   filepath.Base(toFile),
			})
			latestDiffBeginning = i
		} else if i <= latestDiffBeginning+1 {
			// The line where the type of change is declared
			diffs[len(diffs)-1].Change = line
		} else if strings.HasPrefix(line, "---") || strings.HasPrefix(line, "+++") {
			// These are just two useless lines, ignore
		} else if strings.HasPrefix(line, "@@") {
			// Chunk header
			without := strings.Replace(line, "@@", "", 1)
			index := strings.Index(without, "@@") + 2

			var lines []*Line
			finalLine := ""
			if index+3 < len(line) {
				finalLine = line[index+3 : len(line)]
			}

			header := line[3 : index-1]
			parts := strings.Split(header, " ")
			start, _ := strconv.Atoi(strings.Split(parts[0][1:], ",")[0])

			// Calculate the starting line

			line := start + i - latestDiffBeginning - 5
			part := Part{
				Header:        header,
				Start:         start,
				LatestAdded:   line,
				LatestRemoved: line,
				Lines: append(lines, &Line{
					raw:       finalLine,
					OldNumber: line,
					NewNumber: line,
					Type:      LineUnchanged,
				}),
			}

			latestDiff := diffs[len(diffs)-1]
			latestDiff.Parts = append(latestDiff.Parts, &part)
		} else {
			// Simple diff line
			latestPart := getLatestPart(diffs)
			if latestPart == nil {
				shared.GetLogger().Info("Igoring git diff line", zap.String("line", line))
			} else {
				lineType := LineUnchanged
				if strings.HasPrefix(line, "-") {
					lineType = LineDeletion
					deletions++
					latestPart.LatestRemoved++
				} else if strings.HasPrefix(line, "+") {
					lineType = LineAddition
					insertions++
					latestPart.LatestAdded++
				} else {
					// Increment both linecounts if nothing changes
					latestPart.LatestAdded++
					latestPart.LatestRemoved++
				}

				latestPart.Lines = append(latestPart.Lines, &Line{
					raw:       line,
					OldNumber: latestPart.LatestAdded,
					NewNumber: latestPart.LatestAdded,
					Type:      lineType,
				})
			}
		}
	}

	c.HTML(200, "diff.tmpl", gin.H{
		"username":     username,
		"user":         c.Keys["user"],
		"selecteddiff": true,
		"repo":         _repo.Name,
		"isownrepo":    shared.HasAccess(c, []string{"repo:settings"}, username, _repo.Name),

		"id":          sha,
		"shortid":     sha[:8],
		"author":      commit.Author,
		"name":        commit.Summary(),
		"description": strings.Replace(commit.Message(), commit.Summary()+"\n", "", 1),
		"time":        humanize.Time(commit.Author.When),

		"sha":        sha,
		"diff":       diffs,
		"deletions":  deletions,
		"insertions": insertions,
	})
}
