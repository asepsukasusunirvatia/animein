package utils

import (
	"fmt"
	"math/rand/v2"
	"net/http"

	"codeberg.org/Asep5K/animein/models"
)

func createNewRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("utils.createNewRequest Error: %w", err)
	}
	setHeaders(req)
	return req, nil
}

func Reqwest[T any](url string, queryParams models.Dict) (T, error) {
	var nill T
	req, err := createNewRequest(url)
	if err != nil {
		return nill, err
	}
	setQuery(req, queryParams)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nill, fmt.Errorf("utils.Reqwest Error: %w", err)
	}
	if res.StatusCode != 200 {
		return nill, fmt.Errorf("\r\033[5m\033[91mError: %v\033[0m", res.Status)
	}
	defer res.Body.Close()
	result, err := JsonDecoder[T](res)
	if err != nil {
		return nill, fmt.Errorf("Error: %w", err)
	}
	return result, nil
}

func setHeaders(req *http.Request) {
	auth := models.Dict{
		"Accept":           "application/json, text/plain, */*",
		"Referer":          "https://animeinweb.com",
		"User-Agent":       fmt.Sprintf("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%d.0.0.0 Safari/537.36", 137+rand.IntN(7)),
		"x-proxy-secret":   "animein-secure-proxy-key-123",
		"Proxy-Connection": "keep-alive",
	}
	for key, value := range auth {
		req.Header.Set(key, value)
	}
}

func setQuery(req *http.Request, params models.Dict) {
	if len(params) == 0 {
		return
	}
	query := req.URL.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	req.URL.RawQuery = query.Encode()
}

// vim: ft=go
