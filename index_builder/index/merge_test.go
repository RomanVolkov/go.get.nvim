package index_test

import "main/index"
import "testing"

func TestMerge(t *testing.T) {

	set1 := map[string]bool{"val1": true, "val2": true, "val3": true}
	set2 := map[string]bool{"val1": true, "val4": true}

	merged := index.Merge(set1, set2)
	if len(merged) != 4 {
		t.Error("incorrect merging")
	}

}
