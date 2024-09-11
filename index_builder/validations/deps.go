package validations

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const depsDevURL = "https://api.deps.dev"
const getPackageURL = "/v3/systems/go/packages/"

func IsPackageURLValid(packageURL string) (bool, error) {
	client := http.Client{}

	apiURL := fmt.Sprintf("%v%v%v", depsDevURL, getPackageURL, url.PathEscape(packageURL))
	resp, err := client.Get(apiURL)
	if err != nil {
		return true, err
	}

	v := struct {
		Package struct {
			Name string `json:"name"`
		} `json:"packageKey"`
	}{}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(body, &v)

	return v.Package.Name == packageURL, nil
}

func ValidatePackageURLs(urls []string) map[string]bool {
	res := make(map[string]bool)
	var wg sync.WaitGroup
	lock := sync.Mutex{}
	for i, url := range urls {
		wg.Add(1)
		go func(url string, index int, res *map[string]bool) {
			isValid, err := IsPackageURLValid(url)
			if err != nil {
				fmt.Printf("Error validating package %v - error: %v\n", url, err)
			}
			lock.Lock()
			(*res)[url] = isValid
			lock.Unlock()
			wg.Done()

		}(url, i, &res)
	}
	wg.Wait()

	return res
}

func CleanupInvalidPackageURLs(uniqueURLs *map[string]bool) {
	batchSize := 1000

	uniquePackagesSlice := make([]string, 0)
	for v := range *uniqueURLs {
		uniquePackagesSlice = append(uniquePackagesSlice, v)
	}

	for i := 0; i < len(uniquePackagesSlice); i = i + batchSize {
		end := i + batchSize
		if end > len(uniquePackagesSlice) {
			end = len(uniquePackagesSlice)
		}

		validatedPackages := ValidatePackageURLs(uniquePackagesSlice[i:end])
		for url, isValid := range validatedPackages {
			if !isValid {
				fmt.Printf("Removing invalid package url: %v\n", url)
				delete(*uniqueURLs, url)
			}
		}
		time.Sleep(time.Second * 1)
	}
}
