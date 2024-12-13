package utils

import (
	"testing"
)

func Test_GetLatestTags(t *testing.T) {
	tags, err := GetLatestTags()
	if err != nil {
		t.Fatal(err)
		return
	}
	if len(tags) == 0 {
		t.Fatal("tags is empty")
	}
	t.Logf("tags: %v", tags)
}
