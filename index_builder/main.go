package main

import (
	"fmt"
	goindexloader "main/go_index_loader"
	"main/index"
	"time"
)

const indexPath = "../lua/go_get/index.txt"

// TODO:
// 1. validate the URL somehow?: 1) url can be used; 2) discard all forks (gh)
// 1.5 Publish on github
// 2. Run the program once in a week on Github Actions: run & commit into repo
// 3. Commit index into git repo (lfs?)
// 4. Maybe copy it inside plugin folder?

func main() {
	loadedIndex, indexTimestamp, _ := index.LoadIndex(indexPath)
	count := len(loadedIndex)
	fmt.Println(fmt.Sprintf("Index size: %v", count))
	end := time.Now()

	fmt.Println(fmt.Sprintf("start: %v - end: %v", indexTimestamp.Format(time.RFC3339Nano), end.Format(time.RFC3339Nano)))
	uniqueURLs := goindexloader.GetUniqueURLs(indexTimestamp, end, time.Hour*2)
	updatedIndex := index.Merge(loadedIndex, uniqueURLs)
	fmt.Printf("New size: %v ", len(updatedIndex))
	fmt.Printf("diff: %v\n", len(updatedIndex)-count)

	index.StoreIndex(indexPath, end, updatedIndex)
}
