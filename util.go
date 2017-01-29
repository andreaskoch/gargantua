package main

import (
	"io/ioutil"
	"net/http"
)

func readURL(url string) ([]byte, error) {
	resp, fetchErr := http.Get(url)
	if fetchErr != nil {
		return nil, fetchErr
	}

	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	return body, nil
}
