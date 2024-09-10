package validations

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

const getRepoURL = "https://api.github.com/repos/"

func getOwnerAndRepo(u string) (string, string) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", ""
	}

	parts := strings.Split(parsed.Path, "/")
	if len(parts) > 2 {
		return parts[1], parts[2]
	}
	return "", ""
}

func GetPackageURL(u string) string {
	parsed, err := url.Parse(u)
	if err != nil {
		return ""
	}

	parts := strings.Split(parsed.Path, "/")
	if len(parts) > 2 {
		return strings.Join(parts[:3], "/")
	}
	return ""
}

// bool - true/false if repo is a fork
// error - any error
func isGithubFork(url string) (bool, error) {
	owner, repo := getOwnerAndRepo(url)
	if len(owner) == 0 {
		return false, errors.New("failed to get owner and repo name from package")
	}

	getURL := getRepoURL + owner + repo
	client := http.Client{}
	req, err := http.NewRequest("GET", getURL, &strings.Reader{})
	if err != nil {
		return false, err
	}
	token := os.Getenv("GITHUB_TOKEN")
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	v := struct {
		IsFork bool `json:"fork"`
	}{}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(body, &v)
	if err != nil {
		return false, err
	}

	return v.IsFork, nil
}

func RemoveForks(uniqueURLs *map[string]bool) {
	// key - repo url; []string - packages
	uniquePackages := map[string][]string{}
	for url := range *uniqueURLs {
		packageURL := GetPackageURL(url)
		if len(packageURL) == 0 {
			continue
		}
		_, exists := uniquePackages[packageURL]
		if !exists {
			uniquePackages[packageURL] = make([]string, 0)
		}
		uniquePackages[packageURL] = append(uniquePackages[packageURL], url)
	}

	var wg sync.WaitGroup
	var lock sync.Mutex = sync.Mutex{}
	uniquePackagesSlice := make([]string, 0)
	for v := range uniquePackages {
		uniquePackagesSlice = append(uniquePackagesSlice, v)
	}

	batchSize := 200
	for i := 0; i < len(uniquePackagesSlice); i = i + batchSize {
		end := i + batchSize
		if end > len(uniquePackagesSlice) {
			end = len(uniquePackagesSlice)
		}

		for _, packageURL := range uniquePackagesSlice[i:end] {
			urls := uniquePackages[packageURL]
			wg.Add(1)
			go func(packageURL string, urls []string, uniqueURLs *map[string]bool, lock *sync.Mutex) {
				if strings.Contains(packageURL, "github.com") {
					isFork, err := isGithubFork(packageURL)
					if err != nil {
						fmt.Println(err)
					} else if isFork {
						fmt.Printf("%v is fork\n", packageURL)
						lock.Lock()
						for _, url := range urls {
							delete(*uniqueURLs, url)
							fmt.Println("removing...", url)
						}
						lock.Unlock()
					}
				}
				wg.Done()
			}(packageURL, urls, uniqueURLs, &lock)
		}
		wg.Wait()
	}
}
