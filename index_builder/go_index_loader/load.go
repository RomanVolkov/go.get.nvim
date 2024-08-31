package goindexloader

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"main/utils"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const indexURL = "https://index.golang.org/index"
const maxConcurentRequests = 10

func getIndex(since time.Time) ([]string, error) {
	result := make([]string, 0)
	url, err := url.Parse(indexURL)

	if err != nil {
		return result, err
	}

	q := url.Query()
	q.Set("since", since.UTC().Format(time.RFC3339Nano))
	url.RawQuery = q.Encode()

	res, err := http.Get(url.String())
	if err != nil {
		return result, err
	}
	if res.StatusCode > 299 {
		return result, errors.New(fmt.Sprintf("Response failed with status code: %d", res.StatusCode))
	}

	scanner := bufio.NewScanner(res.Body)
	defer res.Body.Close()

	for scanner.Scan() {
		var i = struct {
			Path string `json:"Path"`
		}{}
		json.Unmarshal(scanner.Bytes(), &i)
		result = append(result, i.Path)
	}
	if err := scanner.Err(); err != nil {
		return make([]string, 0), err
	}

	return result, nil
}

func loadIndex(since time.Time, wg *sync.WaitGroup, ch chan []string) {
	d, err := getIndex(since)
	if err != nil {
		log.Fatal(err)
		return
	}
	wg.Done()
	ch <- d
}

func GetUniqueURLs(start, end time.Time, duration time.Duration) map[string]bool {
	times := utils.MakeTimeRange(start, end, time.Hour*2)

	var wg sync.WaitGroup
	temp := make(chan []string, len(times))
	sem := make(chan struct{}, maxConcurentRequests)
	for _, t := range times {
		wg.Add(1)
		sem <- struct{}{}
		go func(t time.Time) {
			loadIndex(t, &wg, temp)
			<-sem
		}(t)
	}

	go func() {
		wg.Wait()
		close(temp)
		close(sem)
	}()

	uniqueURLs := map[string]bool{}

	for urls := range temp {
		for _, url := range urls {
			_, exists := uniqueURLs[url]
			if !exists {
				uniqueURLs[url] = true
			}
		}
	}

	return uniqueURLs
}
