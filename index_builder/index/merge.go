package index

func Merge(loaded map[string]bool, updates map[string]bool) (map[string]bool, []string) {
	diff := make([]string, 0)
	for url := range updates {
		_, exists := loaded[url]
		if !exists {
			diff = append(diff, url)
		}
		loaded[url] = true
	}

	return loaded, diff
}
