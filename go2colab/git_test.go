package go2colab

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestCloneRepo(t *testing.T) {
	testUrl := "https://github.com/TimothyStiles/poly/"
	// create a temporary directory
	tmpDataDir, err := ioutil.TempDir("./data", "repo-*")

	if err != nil {
		t.Errorf("Error creating temporary directory: %s", err)
	}
	defer os.RemoveAll(tmpDataDir)
	//create directory

	repo, err := cloneRepo(testUrl, tmpDataDir)
	if err != nil {
		t.Errorf("Error cloning repo: %s", err)
	}

	tagged := tagExists("v0.15.0", repo)

	if tagged == false {
		t.Errorf("Error: v0.15.0 tag does not exist")
	}

	taglist := getSortedListOfTags(repo)
	fmt.Println(taglist)

}
