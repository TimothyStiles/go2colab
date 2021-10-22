package go2colab

import (
	"fmt"
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

	// get go version for docker container build
	repo.GoVersion, err = getGoVersion(gitRepo)
	if err != nil {
		return err
	}

	// retrieve all example tests with keyword "tutorial" in them as Notebooks
	notebooks, err := getTutorials(gitRepo)
	if err != nil {
		return err
	}
	fmt.Println(notebooks)
	// file, _ := json.MarshalIndent(notebook, "", " ")
	// base := filepath.Base(notebook.RepoPath)
	// filename := filepath.Join(tmpDataDir, base)
	// _ = ioutil.WriteFile(filename, file, 0644)

	return nil
}

func getGoVersion(gitRepo *git.Repository) (string, error) {
	var version string
	tree, err := gitRepo.Worktree()
	if err != nil {
		return version, err
	}

	// grep Go Version
	goVersionRegex := regexp.MustCompile("go [\\d].*")
	grepResults, err := grepWorkTree(tree, *goVersionRegex)
	if err != nil {
		return version, err
	}

	// iterate over grep results and get the version when we find the first match for go.mod
	for _, result := range grepResults {
		if result.FileName == "go.mod" {
			version = strings.Split(result.Content, " ")[1]
			break
		}
	}
	return version, nil
}

func getTutorials(gitRepo *git.Repository) ([]Notebook, error) {
	var notebooks []Notebook

	tree, err := gitRepo.Worktree()
	if err != nil {
		return notebooks, err
	}

	// grep Tutorial examples for "Tutorial" keyword.
	// Returns slice of git.GrepResult structs which contain filenames
	// containing our keyword.
	tutorialRegex := regexp.MustCompile("Example_basic")
	grepResults, err := grepWorkTree(tree, *tutorialRegex)
	if err != nil {
		return notebooks, err
	}

	// iterate through all grep results. Filter files that don't have
	// "example" and convert the rest to Notebooks that can be written out.
	for _, result := range grepResults {
		// go2colab only looks at example files.
		if strings.Contains(result.FileName, "example") {

			notebook, err := createNotebook(result.FileName, tree)
			if err != nil {
				return notebooks, err
			}
			notebook.RepoPath = result.FileName

			// append to our notebooks to be returned.
			notebooks = append(notebooks, notebook)
		}
	}

	return notebooks, nil
}

func createNotebook(path string, tree *git.Worktree) (Notebook, error) {

	// intialize notebook and the jupyter "cell" where we'll be storing the example code.
	// notebooks can have many source cells but here we're just using one
	// since we're more concerned with have a repl for MVP.
	var notebook Notebook
	var sourceCell Cell
	sourceCell.CellType = "code" // Need to let jupyter know that this cell should be rendered as code.

	// this is us retrieving the file from the in memory filesystem that represents our git tree.
	src, err := tree.Filesystem.Open(path)
	if err != nil {
		return notebook, err
	}

	// close the file after this function piece of control flow ends
	defer src.Close()

	// read source file contents into bytes
	srcBytes, err := ioutil.ReadAll(src)
	if err != nil {
		return notebook, err
	}
	// convert source bytes to string and add to our code cell.
	srcString := string(srcBytes)
	sourceCell.Source = srcString

	// add the source cell to our notebook.
	// Now notebook can be rendered as jupyter notebook when written out.
	notebook.Cells = append(notebook.Cells, sourceCell)

	return notebook, nil
}
