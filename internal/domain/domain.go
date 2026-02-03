package domain

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type Config struct {
	URLs        map[string]bool
	Seen        map[string]bool
	Client      *http.Client
	InputFile   *os.File
	InputWriter *bufio.Writer
}

var endpoints = []string{
	"/careers",
	"/jobs",
}

func (c *Config) LoadSeenURLs() (map[string]bool, map[string]bool, error) {
	seen := map[string]bool{}
	urls := map[string]bool{}

	scanner := bufio.NewScanner(c.InputFile)

	for scanner.Scan() {
		rawURL := scanner.Text()
		urls[rawURL] = true

		u, err := url.Parse(rawURL)
		if err != nil {
			fmt.Println("failed to parse ", scanner.Text())
			continue
		}

		for _, endpoint := range endpoints {
			key := fmt.Sprintf("https://%s%s", u.Host, endpoint)
			seen[key] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return seen, urls, err
	}

	return seen, urls, nil
}

func (c *Config) ValidateURL(urlStr string) error {
	u, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	for _, endpoint := range endpoints {
		searchURL := fmt.Sprintf("https://%s%s", u.Host, endpoint)

		_, ok := c.Seen[searchURL]
		if ok {
			break
		}

		c.Seen[searchURL] = true
		status, respURL, err := c.TestURL(searchURL)
		if err != nil {
			fmt.Println("error fetching ", searchURL, err.Error())
			continue
		}

		fmt.Println(status, respURL)

		if status == 406 {
			return errors.New("GOT 406")
		}

		if status == http.StatusOK {
			_, v := c.URLs[respURL]
			if v {
				break
			}

			c.URLs[respURL] = true
			_, err := c.InputWriter.WriteString(respURL + "\n")
			if err != nil {
				fmt.Println("error writing to out ", err.Error())
				continue
			}
			break
		}

	}

	return nil
}

func (c *Config) TestURL(url string) (int, string, error) {
	req, err := http.NewRequest(
		"GET",
		url,
		bytes.NewBuffer([]byte("")),
	)
	if err != nil {
		return 0, "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := c.Client.Do(req)
	if err != nil {
		return 0, "", err
	}

	defer resp.Body.Close()

	return resp.StatusCode, resp.Request.URL.String(), nil
}
