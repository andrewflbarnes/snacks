// Package http provides HTTP request building helpers and related enums
package http

import (
	"fmt"
	"strings"
)

// RequestBuilder holds various HTTP request related metadata and can be used
// to generate partial HTTP requests
type RequestBuilder struct {
	Verb     verb
	Headers  map[string]string
	Body     string
	Endpoint string
	Proto    Protocol
}

// Build returns a string containing a partial HTTP request including the required header
// and, if the verb permits, an empty line followed by the beginning of the HTTP body. This
// may be used as the beginning of RUDY style attacks.
func (b RequestBuilder) Build() string {
	payload := b.BuildHead()

	// Empty line after headers
	payload += "\n"

	// Add body
	if b.Verb.hasBody() && len(b.Body) > 0 {
		payload += b.Body
	}

	return payload
}

// BuildBytes returns a byte slice containing a partial HTTP request as per Build
func (b RequestBuilder) BuildBytes() []byte {
	return []byte(b.Build())
}

// BuildHead returns a string containing a partial HTTP request including the required headers.
// This may be used as the beginning of Slow Loris style attacks.
func (b RequestBuilder) BuildHead() string {
	var payload string

	// Add request line
	payload = fmt.Sprintf("%s %s %s\n", b.Verb, b.Endpoint, b.Proto)

	// Add headers
	for header, value := range b.Headers {
		payload += fmt.Sprintf("%s: %s\n", header, value)
	}
	// Add content length header if the request should have a body and the body was provided.
	// In the case no body is provided it is expected the caller will provide the content length
	// header.
	if b.Verb.hasBody() {
		include := true
		for key := range b.Headers {
			if strings.ToLower(key) == "content-type" {
				include = false
			}
		}
		if include {
			payload += fmt.Sprintf("Content-Length: %d\n", len(b.Body))
		}
	}

	// Custom snacks headers
	payload += "User-Agent: snacks\n"

	return payload
}

// BuildHeadBytes returns a byte slice containing a partial HTTP request as per BuildHead
func (b RequestBuilder) BuildHeadBytes() []byte {
	return []byte(b.BuildHead())
}

func (b RequestBuilder) String() string {
	return b.Build()
}
