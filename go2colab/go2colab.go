package go2colab

import (
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

func Go2Colab(urlString string) error {
	var repo Repo
	tmpDataDir, err := ioutil.TempDir("../quarantine", "repo-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDataDir)

	// Save repo url string
	repo.Url = urlString

	repo, err = getRepoUrlMeta(repo)
	if err != nil {
		return err
	}

	gitRepo, err := cloneRepo(urlString, tmpDataDir)
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
	return nil
}

func getRepoUrlMeta(repo Repo) (Repo, error) {
	var ErrInvalidURL error = errors.New("invalid URL")

	// Parse the url
	urlStruct, err := url.Parse(repo.Url)

	// Check for errors
	if err != nil {
		return repo, err
	}
	// Check if the url is valid
	if !urlStruct.IsAbs() {
		return repo, ErrInvalidURL
	}

	// Get the repo host
	repo.Host = urlStruct.Host

	// Get the repo's url path
	repoPathString := urlStruct.Path

	// split the path into an array
	repoPathStringTrimmed := strings.Trim(repoPathString, "/")
	repoPath := strings.Split(repoPathStringTrimmed, "/")

	// check if the repo path is long enough
	pathLength := len(repoPath)
	if pathLength < 2 {
		return repo, ErrInvalidURL
	}

	// Get the repo owner
	repo.Owner = string(repoPath[0])

	// Get the repo name
	repo.Name = string(repoPath[1])

	// Get the repo's commit hash
	if pathLength > 3 {
		repo.Head = string(repoPath[3])
	} else {
		repo.UseLatestReleaseTag = true
	}
	return repo, nil
}
