package go2colab

import (
	"errors"
	"io"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

func cloneRepoMemory(url string) (*git.Repository, error) {
	fs := memfs.New()
	memoryStore := memory.NewStorage()
	openRepo, err := git.Clone(memoryStore, fs, &git.CloneOptions{
		URL: url,
	})

	return openRepo, err
}

func tagExists(tag string, r *git.Repository) bool {

	tags, err := r.TagObjects()
	if err != nil {
		return false
	}
	for {
		t, err := tags.Next()
		if err != nil {
			break
		}
		if t.Name == tag {
			return true
		}

	}
	return false
}

func getSortedListOfTags(r *git.Repository) []TagInfo {
	tags, err := r.TagObjects()
	if err != nil {
		return nil
	}
	var tagsList []TagInfo
	for {
		tag, err := tags.Next()
		if err != nil {
			break
		}
		var tagInfo TagInfo
		tagInfo.Date = tag.Tagger.When
		tagInfo.Name = tag.Name
		tagInfo.Message = tag.Message
		tagInfo.Hash = tag.Hash.String()
		tagInfo.CommitHash = tag.Target.String()
		tagsList = append(tagsList, tagInfo)
	}
	sortedTagsList := sortTagsByDate(tagsList)
	return sortedTagsList
}

func sortTagsByDate(tags []TagInfo) []TagInfo {
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Date.After(tags[j].Date)
	})
	return tags
}

type TagInfo struct {
	Name       string
	Date       time.Time
	Hash       string
	CommitHash string
	Message    string
}

func checkoutCommit(r *git.Repository, commitHash string) error {
	plumbHash := plumbing.NewHash(commitHash)
	commit, err := r.CommitObject(plumbHash)
	if err != nil {
		return err
	}
	workTree, err := r.Worktree()
	if err != nil {
		return err
	}
	err = workTree.Checkout(&git.CheckoutOptions{
		Hash: commit.Hash,
	})
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
