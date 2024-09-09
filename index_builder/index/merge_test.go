package index_test

import (
	"main/index"
	"testing"
)

func TestMerge(t *testing.T) {
	set1 := map[string]bool{"val1": true, "val2": true, "val3": true}
	set2 := map[string]bool{"val1": true, "val4": true}

	merged, diff := index.Merge(set1, set2)
	if len(merged) != 4 {
		t.Error("incorrect merging")
	}
	if len(diff) != 1 {
		t.Error("incorrect merging")
	}

}
