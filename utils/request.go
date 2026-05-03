package utils

import (
	"fmt"
	"math/rand/v2"
	"net/http"
)

func Request(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("utils.Request failed: %w", err)
	}
	setHeaders(req)
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
		return nil, fmt.Errorf("utils.SearchRequest failed: %w", err)
	}
	setHeaders(req)

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
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("utils.SearchRequest failed: %d (%s)", res.StatusCode, res.Status)
	}
	return res, nil
}

func setHeaders(req *http.Request) {
	userAgent := fmt.Sprintf("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%d.0.0.0 Safari/537.36", 137+rand.IntN(7))
	req.Header.Set("Proxy-Connection", "keep-alive")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("x-proxy-secret", "animein-secure-proxy-key-123")
	req.Header.Set("Referer", "https://animeinweb.com/")
}

// vim: ft=go
