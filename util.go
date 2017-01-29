package main

import (
	"io/ioutil"
	"net/http"
	"time"
)

type Response struct {
	body       []byte
	statusCode int
	duration   time.Duration
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

func readURL(url string) (Response, error) {
	startTime := time.Now()
	resp, fetchErr := http.Get(url)
	if fetchErr != nil {
		return Response{}, fetchErr
	}

	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return Response{}, readErr
	}

	duration := time.Since(startTime)

	return Response{
		body:       body,
		statusCode: resp.StatusCode,
		duration:   duration,
	}, nil
}
