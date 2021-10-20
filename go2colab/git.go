package go2colab

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func cloneRepo(url string, path string) (*git.Repository, error) {
	openRepo, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
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
		fmt.Println(t.Name)
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
