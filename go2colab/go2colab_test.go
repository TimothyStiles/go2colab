package go2colab

import "testing"

func TestGo2Colab(t *testing.T) {
	testUrl := "https://github.com/TimothyStiles/poly/"

	// Test the go2colab function
	err := Go2Colab(testUrl)
	if err != nil {
		t.Errorf("go2colab failed with error: %s", err)
	}
}
