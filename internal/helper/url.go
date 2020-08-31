package helper

import (
	"net/url"
	"strings"
)

func ParseUrl(urlString string) *url.URL {
	if len(urlString) == 0 {
		defaultURL, _ := url.Parse("http://localhost:80/")
		return defaultURL
	}

	// Required by RFC3986 parsing
	if !strings.HasPrefix(urlString, "http://") && !strings.HasPrefix(urlString, "https://") {
		urlString = "http://" + urlString
	}

	url, err := url.Parse(urlString)
	if err != nil {
		logger.Fatal(err)
	}

	fullURL := ""
	if len(url.Hostname()) == 0 {
		fullURL += "localhost"
	} else {
		fullURL += url.Hostname()
	}
	fullURL += ":"
	if len(url.Port()) == 0 {
		fullURL += "80"
	} else {
		fullURL += url.Port()
	}

	url, err = url.Parse(urlString)
	if err != nil {
		logger.Fatal(err)
	}

	return url
}
