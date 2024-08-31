package index

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

// store the index as lua table
func StoreIndex(path string, timestamp time.Time, urls map[string]bool) error {
	if len(urls) == 0 {
		return errors.New("Empty index")
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(fmt.Sprintf("%v\n", timestamp.Format(time.RFC3339Nano)))

	for url := range urls {
		f.WriteString(fmt.Sprintf("%v\n", url))
	}
	return nil
}

// load lua table and parse it into set of urls
func LoadIndex(path string) (map[string]bool, time.Time, error) {
	f, err := os.Open(path)
	if err != nil {
		return map[string]bool{}, time.Time{}, err
	}
	defer f.Close()
	index := map[string]bool{}

	scanner := bufio.NewScanner(f)
	lineIndex := 0
	timestamp := time.Time{}
	for scanner.Scan() {
		line := string(scanner.Bytes())
		if lineIndex == 0 {
			// parse time
			parsedTimed, err := time.Parse(time.RFC3339Nano, line)
			if err != nil {
				return map[string]bool{}, time.Time{}, err
			}
			timestamp = parsedTimed
		} else {
			// parse the line
			s := strings.Trim(line, "\n")
			index[s] = true
		}

		lineIndex++
	}

	return index, timestamp, nil
}
