package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const depsDevURL = "https://api.deps.dev"
const getProjectURL = "/v3alpha/projects/"

// so here I want to get and then show
// 1. homepage
// 2. license
// 3. description
type PackageInfo struct {
	ProjectKey struct {
		Id string `json:"id"`
	} `json:"projectKey"`
	Homepage    string `json:"homepage"`
	License     string `json:"license"`
	Description string `json:"description"`
}

// if url is missing inside result map - no information.
func GetPackageInfos(urls []string) map[string]PackageInfo {

	res := map[string]PackageInfo{}

	for _, url := range urls {

		info, err := GetPackageInfo(url)
		if err != nil {
			fmt.Println(err)
			continue
		}

		res[url] = info

	}
	// no slicing for now, I just want to test
	// will add it next time
	// TODO: slicing
	return res
}

func GetPackageInfo(packageURL string) (PackageInfo, error) {
	client := http.Client{}

	apiURL := fmt.Sprintf("%v%v%v", depsDevURL, getProjectURL, url.PathEscape(packageURL))
	resp, err := client.Get(apiURL)
	if err != nil {
		return PackageInfo{}, err
	}
	defer resp.Body.Close()

	v := PackageInfo{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PackageInfo{}, err
	}
	if err = json.Unmarshal(body, &v); err != nil {
		return PackageInfo{}, err
	}
	// removing ; sign as it's used for separation within the index
	v.Description = strings.ReplaceAll(v.Description, ";", " ")
	return v, nil
}
