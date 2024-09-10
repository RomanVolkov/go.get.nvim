package validations

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

// 1. create an empty random project
// 2. go mod init main
// 3. go get <url>
// 4. return the result
// 5. cleanup go cache

func generateGUID() string {
	b := make([]byte, 16)
	rand.Read(b)

	uuid := fmt.Sprintf("%x%x%x%x%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:],
	)
	return uuid
}

func createTempDir() (string, error) {
	tempPath := os.TempDir() + "/" + generateGUID()

	err := os.MkdirAll(tempPath, 0755)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to create directory: %v\n", err))
	}

	return tempPath, nil
}

func createEmptyProject(path string) error {
	cmd := exec.Command("go", "mod", "init", "main")
	cmd.Dir = path
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func cleanupCache() error {
	fmt.Println("Cleaning...")
	cmd := exec.Command("go", "clean", "-modcache")
	fmt.Println(cmd)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func ValidatePackage(url string) (bool, error) {
	tempDir, err := createTempDir()
	if err != nil {
		return false, err
	}

	if err = createEmptyProject(tempDir); err != nil {
		return false, err
	}
	defer os.RemoveAll(tempDir)

	cmd := exec.Command("go", "get", "-u", url)
	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))

	if err != nil {
		return false, err
	}

	return true, nil
}

func ValidatePackages(urls []string) map[string]bool {
	res := make(map[string]bool)
	var wg sync.WaitGroup
	lock := sync.Mutex{}
	for i, url := range urls {
		wg.Add(1)
		go func(url string, index int, res *map[string]bool) {
			isValid, err := ValidatePackage(url)
			if err != nil {
				fmt.Printf("Error validating package %v - error: %v\n", url, err)
			}
			lock.Lock()
			defer lock.Unlock()
			defer wg.Done()
			(*res)[url] = isValid

		}(url, i, &res)
	}
	wg.Wait()

	if err := cleanupCache(); err != nil {
		fmt.Printf("Error cleaning go modcache: %v\n", err)
	}

	return res
}

func CleanupInvalidPackages(uniqueURLs *map[string]bool) []string {
	batchSize := 100
	removedURLs := make([]string, 0)

	uniquePackagesSlice := make([]string, 0)
	for v := range *uniqueURLs {
		uniquePackagesSlice = append(uniquePackagesSlice, v)
	}

	for i := 0; i < len(uniquePackagesSlice); i = i + batchSize {
		end := i + batchSize
		if end > len(uniquePackagesSlice) {
			end = len(uniquePackagesSlice)
		}

		validatedPackages := ValidatePackages(uniquePackagesSlice[i:end])
		for url, isValid := range validatedPackages {
			if !isValid {
				fmt.Printf("Removing invalid package url: %v\n", url)
				delete(*uniqueURLs, url)
				removedURLs = append(removedURLs, url)
			}
		}
		time.Sleep(time.Second)
	}
	return removedURLs
}
