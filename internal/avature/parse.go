package avature

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/dhaef/job-scraper/internal/dom"
	"github.com/dhaef/job-scraper/internal/job"
	"github.com/dhaef/job-scraper/internal/request"
	"golang.org/x/net/html"
)

func GetPaginationLinks(doc *html.Node) []*html.Node {
	return dom.FindElementsByClasses(doc, []string{"paginationNextLink"})
}

func FetchJobsPage(url string) (*html.Node, error) {
	responseBody, err := request.Fetch("GET", url, []byte(""), map[string]string{
		"User-Agent": "Mozilla/5.0",
	})
	if err != nil {
		return nil, fmt.Errorf("error making request: %s", err.Error())
	}

	bodyAsStr := string(responseBody)
	if bodyAsStr == "" {
		return nil, fmt.Errorf("no body found for url: %s", url)
	}

	doc, err := html.Parse(bytes.NewReader(responseBody))
	if err != nil {
		return nil, fmt.Errorf("parsing html: %s", err.Error())
	}

	return doc, nil
}

func GetJobs(doc *html.Node) ([]*html.Node, int, error) {
	jobs := dom.FindElementsByClasses(doc, []string{"article__header__text__title", "list__item__text__title"})

	return jobs, len(jobs), nil
}

func FetchJobDetailsPage(url string) (*html.Node, error) {
	jobDetailsResponse, err := request.Fetch("GET", url, []byte(""), map[string]string{
		"User-Agent": "Mozilla/5.0",
	})
	if err != nil {
		return nil, fmt.Errorf("error making details request: %s, url: %s", err.Error(), url)
	}

	detailsRespAsStr := string(jobDetailsResponse)
	if detailsRespAsStr == "" {
		return nil, fmt.Errorf("no body found for url: %s", url)
	}

	detailsDoc, err := html.Parse(bytes.NewReader(jobDetailsResponse))
	if err != nil {
		return nil, fmt.Errorf("parsing html: %v", err)
	}

	return detailsDoc, nil
}

func GetLinkAndTitle(doc *html.Node) (*job.Job, error) {
	linkElements := dom.FindElementsByTagName(doc, "a")
	if len(linkElements) == 0 {
		return nil, errors.New("no link elements found")
	}

	linkElement := linkElements[0]

	link, ok := dom.GetHref(linkElement)
	if !ok {
		return nil, errors.New("no link found from link element")
	}
	job := job.NewJob()
	job.SetTitle(strings.ToLower(dom.GetTextContent(linkElement)))
	job.SetLink(link)

	return job, nil
}

func GetDescription(doc *html.Node, j *job.Job) {
	details := dom.FindElementsByClasses(doc, []string{"article--details"})
	j.SetDescription(details)
}

func GetMetadata(doc *html.Node, j *job.Job, url string) {
	metadataNodes := dom.FindElementsByClasses(doc, []string{"article__content__view__field"})

	for _, item := range metadataNodes {
		labelDocs := dom.FindElementsByClasses(item, []string{"article__content__view__field__label"})
		valueDocs := dom.FindElementsByClasses(item, []string{"article__content__view__field__value"})

		var label string
		var value string

		if len(labelDocs) == 0 && len(valueDocs) == 0 {
			continue
		} else if len(labelDocs) >= 1 && len(valueDocs) == 0 {
			parts := strings.Split(strings.TrimSpace(dom.GetTextContent(labelDocs[0])), ":")
			if len(parts) == 2 {
				label = parts[0]
				value = parts[1]
			}
		} else if len(labelDocs) == 0 && len(valueDocs) >= 1 {
			parts := strings.Split(strings.TrimSpace(dom.GetTextContent(valueDocs[0])), ":")
			if len(parts) == 2 {
				label = parts[0]
				value = parts[1]
			} else if len(parts) == 1 {
				l, ok := dom.FindLabelByClass(item, job.LabelMap)
				if ok {
					label = l
					value = strings.TrimSpace(dom.GetTextContent(valueDocs[0]))
				}
			}
		} else {
			label = strings.TrimSpace(dom.GetTextContent(labelDocs[0]))
			value = strings.TrimSpace(dom.GetTextContent(valueDocs[0]))
		}

		normalizedLabel, ok := job.NormalizeMetadataLabel(strings.ToLower(label))
		if ok {
			j.SetMetadata(normalizedLabel, strings.ToLower(value))
		}
	}

	company, ok := j.Metadata["company"]
	if !ok || company == "" {
		j.Metadata["company"] = ExtractUniqueHostPart(url)
	}
}

var excludedURLParts = map[string]bool{
	"www":         true,
	"apply":       true,
	"jobs":        true,
	"careers":     true,
	"net":         true,
	"com":         true,
	"avature":     true,
	"org":         true,
	"recruitment": true,
}

func ExtractUniqueHostPart(URL string) string {
	u, err := url.Parse(URL)
	if err != nil {
		return ""
	}

	host := u.Hostname()

	for part := range strings.SplitSeq(host, ".") {
		if !excludedURLParts[strings.ToLower(part)] {
			return part
		}
	}

	return ""
}
