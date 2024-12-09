package model

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func MakeRequest(method, url, bearer, data string) (*http.Response, error) {
	var req *http.Request
	var err error
	if data != "" {
		req, err = http.NewRequest(method, url, bytes.NewBuffer([]byte(data)))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("error creating the request: %v", err)
	}

	if data != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Add("Authorization", "Bearer "+bearer)

	client := &http.Client{}
	return client.Do(req)
}

func ReadResponseBody(body io.ReadCloser) ([]byte, error) {
	return io.ReadAll(body)
}
