package index

import (
	"bufio"
	"errors"
	"fmt"
	"main/utils"
	"os"
	"strings"
	"sync"
	"time"
)

type IndexItem struct {
	URL         string
	Homepage    string
	License     string
	Description string
}

// store the index as lua table
// url:IndexItem
func StoreIndex(path string, timestamp time.Time, urls map[string]IndexItem) error {
	if len(urls) == 0 {
		return errors.New("Empty index")
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(fmt.Sprintf("%v\n", timestamp.Format(time.RFC3339Nano)))

	for url, info := range urls {
		// I am going to use ; separator (CSV-style) to store the data
		f.WriteString(fmt.Sprintf("%v;%v;%v;%v\n", url, info.License, info.Homepage, info.Description))
	}
	return nil
}

// load lua table and parse it into set of urls
// url:IndexItem
func LoadIndex(path string) (map[string]IndexItem, time.Time, error) {
	f, err := os.Open(path)
	if err != nil {
		return map[string]IndexItem{}, time.Time{}, err
	}
	defer f.Close()
	index := map[string]IndexItem{}

	scanner := bufio.NewScanner(f)
	lineIndex := 0
	timestamp := time.Time{}
	for scanner.Scan() {
		line := string(scanner.Bytes())
		if lineIndex == 0 {
			// parse time
			parsedTimed, err := time.Parse(time.RFC3339Nano, line)
			if err != nil {
				return map[string]IndexItem{}, time.Time{}, err
			}
			timestamp = parsedTimed
		} else {
			// parse the line
			s := strings.Trim(line, "\n")

			values := strings.Split(s, ";")
			if len(values) != 4 {
				fmt.Println("not enough values to parse for ", s)
				continue
			}
			// now we need to split the line with ";" and then make an IndexItem
			url := values[0]

			index[url] = IndexItem{
				URL:         url,
				License:     values[1],
				Homepage:    values[2],
				Description: values[3],
			}
		}

		lineIndex++
	}

	return index, timestamp, nil
}

func UpdateIndex(indexData *map[string]IndexItem, uniqueURLs map[string]bool) {
	var wg sync.WaitGroup
	var lock sync.Mutex = sync.Mutex{}
	slice := make([]string, 0)
	for v := range uniqueURLs {
		slice = append(slice, v)
	}

	batchSize := 20
	for i := 0; i < len(slice); i = i + batchSize {
		end := i + batchSize
		if end > len(slice) {
			end = len(slice)
		}

		for _, url := range slice[i:end] {
			wg.Add(1)
			go func(url string, indexData *map[string]IndexItem, lock *sync.Mutex) {
				packageInfo, err := utils.GetPackageInfo(url)
				defer wg.Done()
				if err != nil {
					fmt.Printf("Error to get packageInfo for %v - %v\n", url, err)
				}
				lock.Lock()
				defer lock.Unlock()
				(*indexData)[url] = IndexItem{
					URL:         url,
					Homepage:    packageInfo.Homepage,
					License:     packageInfo.License,
					Description: packageInfo.Description,
				}
			}(url, indexData, &lock)
		}
		wg.Wait()
	}
}
