package go2colab

import (
	"testing"
)

func TestCloneRepoMemory(t *testing.T) {
	testUrl := "https://github.com/TimothyStiles/poly/"

	repo, err := cloneRepoMemory(testUrl)
	if err != nil {
		t.Errorf("Error cloning repo: %s", err)
	}

	tagged := tagExists("v0.15.0", repo)

	if tagged == false {
		t.Errorf("Error: v0.15.0 tag does not exist")
	}
}
