package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Response struct {
	body        []byte
	statusCode  int
	duration    time.Duration
	contentType string
}

func (response *Response) Body() []byte {
	return response.body
}

func (response *Response) Size() int {
	return len(response.body)
}

func (response *Response) StatusCode() int {
	return response.statusCode
}

func (response *Response) Duration() time.Duration {
	return response.duration
}

func (response *Response) ContentType() string {
	return response.contentType
}

func (response *Response) IsHTML() bool {
	return strings.HasPrefix(response.contentType, "text/html")
}

func readURL(url url.URL) (Response, error) {
	startTime := time.Now()
	resp, fetchErr := http.Get(url.String())
	if fetchErr != nil {
		return Response{}, fetchErr
	}

	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return Response{}, readErr
	}

	duration := time.Since(startTime)

	// content type
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(body)
	}

	return Response{
		body:        body,
		statusCode:  resp.StatusCode,
		duration:    duration,
		contentType: contentType,
	}, nil
}
