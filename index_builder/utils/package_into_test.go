package utils_test

import (
	"fmt"
	"main/utils"
	"testing"
)

// Okay, it works
// so now I can get all the package into data (if it exists) and add it into the index
func TestGetPackageInfo(t *testing.T) {
	testURL := "github.com/go-chi/chi"

	info, err := utils.GetPackageInfo(testURL)
	if err != nil {
		t.Error(err)
	}

	if info.ProjectKey.Id != testURL {
		fmt.Println(info.ProjectKey.Id)
		t.Error("incorrect ProjectKey")
	}

}
