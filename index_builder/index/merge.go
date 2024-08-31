package index

import "strings"

func Merge(loaded map[string]bool, updates map[string]bool) map[string]bool {
	for url := range updates {
		loaded[strings.ToLower(url)] = true
	}

	return loaded
}
