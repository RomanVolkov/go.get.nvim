package main

import (
	"fmt"
	goindexloader "main/go_index_loader"
	"main/index"
	"main/validations"
	"os"
	"time"

	"github.com/OlyMahmudMugdho/gotenv/gotenv"
)

const indexPath = "../lua/go_get/index.txt"

func main() {
	_, err := os.Stat(".env")
	if !os.IsNotExist(err) {
		fmt.Println("loading env...")
		gotenv.Load()
	}

	loadedIndex, indexTimestamp, _ := index.LoadIndex(indexPath)
	count := len(loadedIndex)
	fmt.Println(fmt.Sprintf("Index size: %v", count))
	end := time.Now()

	end = indexTimestamp.Add(time.Hour * 24 * 31)

	fmt.Println(fmt.Sprintf("start: %v - end: %v", indexTimestamp.Format(time.RFC3339Nano), end.Format(time.RFC3339Nano)))
	uniqueURLs := goindexloader.GetUniqueURLs(indexTimestamp, end, time.Hour*2)
	fmt.Println(fmt.Sprintf("Loaded urls size: %v", len(uniqueURLs)))

	fmt.Println("Cleaning forks...")
	validations.RemoveForks(&uniqueURLs)

	fmt.Println("Cleaning incorrect package urls...")
	countBefore := len(uniqueURLs)
	validations.CleanupInvalidPackageURLs(&uniqueURLs)
	fmt.Println("Removed %v incorrect package urls", countBefore-len(uniqueURLs))

	fmt.Println("Cleaning invalid packages...")
	removedURLs := validations.CleanupInvalidPackages(&uniqueURLs)
	fmt.Println("Removed %v incorrect packages", len(removedURLs))
	fmt.Println(removedURLs)

	updatedIndex, diff := index.Merge(loadedIndex, uniqueURLs)
	fmt.Printf("New size: %v ", len(updatedIndex))
	fmt.Printf("diff: %v\n", len(updatedIndex)-count)
	fmt.Println(diff)

	index.StoreIndex(indexPath, end, updatedIndex)
}
