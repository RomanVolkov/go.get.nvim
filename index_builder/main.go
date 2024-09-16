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

	indexData, indexTimestamp, _ := index.LoadIndex(indexPath)
	count := len(indexData)
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
	fmt.Printf("Removed %v incorrect package urls\n", countBefore-len(uniqueURLs))

	// fmt.Println("Cleaning invalid packages...")
	// removedURLs := validations.CleanupInvalidPackages(&uniqueURLs)
	// fmt.Println("Removed %v incorrect packages", len(removedURLs))
	// fmt.Println(removedURLs)

	index.UpdateIndex(&indexData, uniqueURLs)

	fmt.Printf("New index size: %v; diff: %v\n", len(indexData), len(indexData)-count)

	index.StoreIndex(indexPath, end, indexData)
}
