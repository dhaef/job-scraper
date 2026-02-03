package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dhaef/job-scraper/internal/avature"
	"github.com/dhaef/job-scraper/internal/dom"
)

func main() {
	fmt.Println("Starting server...")

	fileIn, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	fileOut, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	defer fileOut.Close()

	scanner := bufio.NewScanner(fileIn)

	for scanner.Scan() {
		url := scanner.Text()

		err = handleURL(url, fileOut)
		if err != nil {
			fmt.Printf("err handling url: %s, err: %v\n", url, err)
		}

		time.Sleep(500 * time.Millisecond)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err.Error())
	}
}

func handleURL(url string, f *os.File) error {
	finalPage := false

	for {
		fmt.Println("Navigating to ", url)

		doc, err := avature.FetchJobsPage(url)
		if err != nil {
			fmt.Printf("failed to get doc, %s\n", err.Error())
			break
		}

		nextPaginationLinks := avature.GetPaginationLinks(doc)
		if len(nextPaginationLinks) == 0 {
			finalPage = true
		} else {
			nextURL, ok := dom.GetHref(nextPaginationLinks[0])
			if !ok {
				finalPage = true
			} else {
				url = nextURL
			}
		}

		jobPostings, numberOfJobsOnPage, err := avature.GetJobs(doc)
		if err != nil {
			fmt.Printf("failed to fetchJobs, %s\n", err.Error())
			break
		}

		fmt.Println("Found jobs: ", numberOfJobsOnPage)

		if numberOfJobsOnPage == 0 {
			break
		}

		for _, jobPosting := range jobPostings {
			time.Sleep(500 * time.Millisecond)

			job, err := avature.GetLinkAndTitle(jobPosting)
			if err != nil {
				fmt.Printf("getting link and title: %v", err)
				continue
			}

			detailsDoc, err := avature.FetchJobDetailsPage(job.Link)
			if err != nil {
				fmt.Printf("fetching job details: %v", err)
				continue
			}

			avature.GetDescription(detailsDoc, job)
			avature.GetMetadata(detailsDoc, job, url)

			jobJSON, err := job.ToJSON()
			if err != nil {
				fmt.Printf("failed to convert job %s to json, err: %v\n", job.Link, err)
				continue
			}

			if _, err := fmt.Fprintf(f, "%s\n", jobJSON); err != nil {
				fmt.Println("failed to append job, ", err.Error())
			}
		}

		if finalPage {
			break
		}
	}

	return nil
}
