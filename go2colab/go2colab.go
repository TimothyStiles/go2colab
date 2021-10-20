package go2colab

import (
	"io/ioutil"
	"os"
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
