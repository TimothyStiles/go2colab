package go2colab

import (
	"regexp"
	"strings"

	"github.com/go-git/go-billy/v5"
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

	var examples []billy.File
	for _, result := range grepResults {
		if strings.Contains(result.FileName, "example") {
			file, err := tree.Filesystem.Open(result.FileName)
			if err != nil {
				return err
			}
			examples = append(examples, file)
		}
	}
	return nil
}
