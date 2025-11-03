package web

import "net/url"

func IsValidURL(rawURL string) bool {
	parsedURL, err := url.ParseRequestURI(rawURL)

	if err != nil {
		return false
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	if parsedURL.Host == "" {
		return false
	}

	return true
}
