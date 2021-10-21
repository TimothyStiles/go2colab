package go2colab

import (
	"fmt"
	"io"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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
	treeStat, err := tree.Status()
	if err != nil {
		return err
	}
	fmt.Println(root, treeStat)
	return nil
}

func getBlobs(repo *git.Repository) ([]BlobInfo, error) {

	// blob object types: https://github.com/go-git/go-git/blob/4ec1753b4e9324d455d3b55060020ce324e6ced2/plumbing/object.go#L42
	// InvalidObject ObjectType = 0
	// CommitObject  ObjectType = 1
	// TreeObject    ObjectType = 2
	// BlobObject    ObjectType = 3
	// TagObject     ObjectType = 4

	var blobDetails []BlobInfo
	blobs, err := repo.BlobObjects()
	if err != nil {
		return nil, err
	}
	for {
		blob, err := blobs.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		var blobInfo BlobInfo
		blobInfo.Hash = blob.Hash.String()
		blobInfo.ObjectType = blob.Type()
		blobInfo.Size = blob.Size
		blobInfo.BlobReader, err = blob.Reader()
		if err != nil {
			return nil, err
		}
		blobDetails = append(blobDetails, blobInfo)
	}
	return blobDetails, nil
}

type BlobInfo struct {
	BlobReader io.ReadCloser
	Hash       string
	ObjectType plumbing.ObjectType
	Size       int64
}
