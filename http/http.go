package http

import (
	"fmt"
)

type HttpRequestBuilder struct {
	Verb     HttpVerb
	Headers  map[string]string
	Body     string
	Endpoint string
	Proto    HttpProto
}

func (b HttpRequestBuilder) GetPayload() string {
	var payload string

	// Add request line
	payload = fmt.Sprintf("%s %s %s\n", b.Verb, b.Endpoint, b.Proto)

	// Add headers
	for header, value := range b.Headers {
		payload += fmt.Sprintf("%s: %s\n", header, value)
	}
	// Add content length header
	// if b.Verb.hasBody() {
	// 	payload += fmt.Sprintf("Content-Length: %d\n", len(b.Body))
	// }

	// Custom snacks headers
	payload += "User-Agent: snacks\n"

	// Post header empty line
	payload += "\n"

	// Add body
	if b.Verb.hasBody() {
		payload += b.Body
	}

	return payload
}

func (b HttpRequestBuilder) GetPayloadBytes() []byte {
	return []byte(b.GetPayload())
}

func (b HttpRequestBuilder) String() string {
	return b.GetPayload()
}
