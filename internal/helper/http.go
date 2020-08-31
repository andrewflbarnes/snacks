package helper

import (
	"github.com/andrewflbarnes/snacks/pkg/http"
	log "github.com/sirupsen/logrus"
)

var (
	payloadPrefixes = map[http.ContentType][]byte{
		http.ApplicationJSON:               []byte(`{"a":"`),
		http.ApplicationXWWWFormURLEncoded: []byte(`a=`),
	}
)

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
