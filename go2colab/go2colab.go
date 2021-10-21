package go2colab

import (
	"fmt"
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
	root := tree.Filesystem.Root()
	files, err := tree.Filesystem.ReadDir(root)
	if err != nil {
		return err
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	fmt.Println(fileNames)
	treeStat, err := tree.Status()
	if err != nil {
		return err
	}
	fmt.Println(root, treeStat)
	return nil
}
