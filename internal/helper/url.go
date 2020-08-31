package helper

import (
	"net/url"
	"strings"
)

// ParseURL parses the string representation of a URL and returns a url.URL pointer.
// The string should comply with RFC3986 with the following exception - if the scheme
// is not set, or not one of http or https, then it will prepend http:// prior to passimg.
// Aside from the scheme, the following defaults are set if not present: port 80, host
// localhost.
// If the urlString is empty then the default http://localhost:80 is used
func ParseURL(urlString string) *url.URL {
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
