package helper

import (
	"github.com/andrewflbarnes/snacks/pkg/http"
	log "github.com/sirupsen/logrus"
)

var (
	payloadPrefixes = map[http.ContentType][]byte{
		http.ApplicationJSON:               []byte(`{"a":"`),
		http.ApplicationXML:                []byte(`<?xml version="1.0" encoding="UTF-8"?><a>`),
		http.ApplicationXWWWFormURLEncoded: []byte(`a=`),
	}
)

// ToContentType converts a string to an http.ContentType representation. If one
// does not currently exist a new ContentType is created and returned
func ToContentType(contentType string) http.ContentType {
	val := http.ToContentType(contentType)

	if val == http.ContentTypeNotFound {
		logger.WithFields(log.Fields{
			"contentType": contentType,
		}).Warn("Content type not found, creating custom type")
		return http.ContentType(contentType)
	}

	return val
}

// GetPayloadPrefix returns a default payload prefix corresponding to a specific
// ContentType. If no default is found specific to the ContentType then a "a="
// is returned.
func GetPayloadPrefix(media http.ContentType) []byte {
	prefix, ok := payloadPrefixes[media]

	if !ok {
		logger.WithFields(log.Fields{
			"media":   media,
			"default": "a=",
		}).Warn("No payload prefix found, defaulting")
		return []byte("a=")
	}

	return prefix
}
