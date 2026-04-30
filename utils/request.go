package utils

import (
	"fmt"
	"net/http"
)

func Request(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("utils.Request failed: %w", err)
	}
	setHeaders(req)
	// req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
	// req.Header.Set("Referer", "https://animeinweb.com/")
	// req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("utils.Request failed: %w", err)
	}

	return res, nil
}

func SearchRequest(keyWord string) (*http.Response, error) {
	targetURL := "https://animeinweb.com/api/proxy/3/2/explore/movie"

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("utils.Request failed: %w", err)
	}
	setHeaders(req)
	// req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
	// req.Header.Set("Referer", "https://animeinweb.com/")
	// req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// Set Query Parameters
	query := req.URL.Query() // Ambil objek query bawaan URL
	query.Add("page", "0")
	query.Add("sort", "views")
	query.Add("keyword", keyWord)

	req.URL.RawQuery = query.Encode()

	client := &http.Client{}
	res, err := client.Do(req) // Eksekusi dengan objek 'req'

	if err != nil {
		return nil, fmt.Errorf("utils.SearchRequest failed: %w", err)
	}

	return res, nil
}

func setHeaders(req *http.Request) {
    req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
    req.Header.Set("Referer", "https://animeinweb.com/")
    req.Header.Set("Content-Type", "application/json; charset=utf-8")
}

// vim: ft=go
