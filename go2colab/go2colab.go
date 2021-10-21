package go2colab

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
)

func Go2Colab(urlString string) error {
	var repo Repo

	// Save repo url string
	repo.Url = urlString

	repo, err := getRepoUrlMeta(repo)
	if err != nil {
		return err
	}

	gitRepo, err := cloneRepoMemory(urlString)
	if err != nil {
		return err
	}

	if repo.UseLatestReleaseTag {
		repo.ReleaseTag = getSortedListOfTags(gitRepo)[0]
		repo.Head = repo.ReleaseTag.CommitHash
	}

	err = checkoutCommit(gitRepo, repo.Head)
	if err != nil {
		return err
	}

	tree, err := gitRepo.Worktree()
	if err != nil {
		return err
	}

	// grep Go Version
	goVersionRegex := regexp.MustCompile("go [\\d].*")
	grepResults, err := grepWorkTree(tree, *goVersionRegex)
	if err != nil {
		return err
	}

	// iterate over grep results and get the version when we find the first match for go.mod
	for _, result := range grepResults {
		if result.FileName == "go.mod" {
			version := strings.Split(result.Content, " ")[1]
			repo.GoVersion = version
			break
		}
	}

	// grep Tutorial examples
	tutorialRegex := regexp.MustCompile("Example_basic")
	grepResults, err = grepWorkTree(tree, *tutorialRegex)
	if err != nil {
		return err
	}

	//  iterate over example files, check filename, and store them as billy.File slice
	var notebooks []Notebook
	for _, result := range grepResults {
		if strings.Contains(result.FileName, "example") {
			notebook, err := grep2Notebook(tree, result)
			if err != nil {
				return err
			}
			notebooks = append(notebooks, notebook)
		}
	}

	return nil
}

func grep2Notebook(tree *git.Worktree, grepResult git.GrepResult) (Notebook, error) {
	var notebook Notebook
	var sourceCell Cell
	sourceCell.CellType = "code"
	src, err := tree.Filesystem.Open(grepResult.FileName)
	if err != nil {
		return notebook, err
	}
	defer src.Close()

	srcBytes, err := ioutil.ReadAll(src)
	if err != nil {
		return notebook, err
	}
	srcString := string(srcBytes)
	sourceCell.Source = srcString
	notebook.Cells = append(notebook.Cells, sourceCell)
	notebook.RepoPath = grepResult.FileName
	return notebook, nil
}
