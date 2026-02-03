package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dhaef/job-scraper/internal/domain"
)

type Record struct {
	URL string `json:"url"`
}

type Index struct {
	ID  string `json:"id"`
	URL string `json:"cdx-api"`
}

func main() {
	file, err := os.OpenFile("input.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	c := domain.Config{
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
		Seen:        map[string]bool{},
		InputFile:   file,
		InputWriter: bufio.NewWriter(file),
		URLs:        map[string]bool{},
	}

	seen, urls, err := c.LoadSeenURLs()
	if err != nil {
		log.Fatal(err)
	}
	c.Seen = seen
	c.URLs = urls

	indexes := getIndexes()

	for _, index := range indexes {
		fmt.Println(index.ID)

		indexURL := fmt.Sprintf("%s?url=*.avature.net&output=json", index.URL)
		resp, err := c.Client.Get(indexURL)
		if err != nil {
			fmt.Println("error fetching index data", err.Error())
			continue
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			var record Record
			if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
				continue
			}

			err := c.ValidateURL(record.URL)
			if err != nil {
				if err.Error() == "GOT 406" {
					log.Fatalf("got 406 at %s", record.URL)
				}

				fmt.Println("failed to test url", record.URL)
				continue
			}

			time.Sleep(500 * time.Millisecond)
		}

		time.Sleep(2 * time.Second)
	}

	c.InputWriter.Flush()
}

func getIndexes() []Index {
	file, err := os.ReadFile("common-crawl-indexes.json")
	if err != nil {
		log.Fatal(err)
	}

	var indexes []Index
	err = json.Unmarshal(file, &indexes)
	if err != nil {
		log.Fatal(err)
	}

	return indexes
}
