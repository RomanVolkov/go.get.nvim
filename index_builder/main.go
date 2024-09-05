package main

import (
	"fmt"
	goindexloader "main/go_index_loader"
	"main/index"
	"main/validations"
	"time"
)

const indexPath = "../lua/go_get/index.txt"

func main() {
	loadedIndex, indexTimestamp, _ := index.LoadIndex(indexPath)
	count := len(loadedIndex)
	fmt.Println(fmt.Sprintf("Index size: %v", count))
	end := time.Now()

	end = indexTimestamp.Add(time.Hour * 24 * 7)
	//
	// fmt.Println(fmt.Sprintf("start: %v - end: %v", indexTimestamp.Format(time.RFC3339Nano), end.Format(time.RFC3339Nano)))
	//
	uniqueURLs := goindexloader.GetUniqueURLs(indexTimestamp, end, time.Hour*2)
	fmt.Println(fmt.Sprintf("Loaded urls size: %v", count))

	validations.RemoveForks(&uniqueURLs)
	// TODO: validate diff packages with go get - with index update
	// TODO: validate whole index? - separate gh action
	// 1. get url
	// 2. create a tmp dir for 1 url and init go mod
	// 3. try to install the package with go get -u
	// 4. get the output - false if failed to install
	// 5. remove tmp dir - but are they stored there?? Cleanup local checkouts after each batch somehow?

	updatedIndex, diff := index.Merge(loadedIndex, uniqueURLs)
	fmt.Printf("New size: %v ", len(updatedIndex))
	fmt.Printf("diff: %v\n", len(updatedIndex)-count)
	fmt.Println(diff)

	index.StoreIndex(indexPath, end, updatedIndex)
}
