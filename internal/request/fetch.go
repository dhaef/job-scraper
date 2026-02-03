package request

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

func Fetch(method string, url string, body []byte, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(
		method,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		log.Fatalf("building request: %v", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("making request: %v", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("reading body: %v", err)
	}

	if resp.StatusCode == 406 {
		log.Fatalf("failed with 406: %s", url)
	}

	if resp.StatusCode != http.StatusOK {
		return respBody, fmt.Errorf("request failed with statusCode: %d", resp.StatusCode)
	}

	return respBody, nil
}
