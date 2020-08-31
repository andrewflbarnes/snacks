// Package http provides HTTP request building helpers and related enums
package http

import (
	"fmt"
	"strings"
)

type HttpRequestBuilder struct {
	Verb     verb
	Headers  map[string]string
	Body     string
	Endpoint string
	Proto    HttpProto
}

func (b HttpRequestBuilder) Build() string {
	payload := b.BuildHead()

	// Empty line after headers
	payload += "\n"

	// Add body
	if b.Verb.hasBody() && len(b.Body) > 0 {
		payload += b.Body
	}

	return payload
}

func (b HttpRequestBuilder) BuildBytes() []byte {
	return []byte(b.Build())
}

func (b HttpRequestBuilder) BuildHead() string {
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
		for key, _ := range b.Headers {
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

func (b HttpRequestBuilder) BuildHeadBytes() []byte {
	return []byte(b.BuildHead())
}

func (b HttpRequestBuilder) String() string {
	return b.Build()
}
