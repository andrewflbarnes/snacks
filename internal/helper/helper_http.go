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
	media := http.ToContentType(contentType)
	if media == http.ContentTypeNotFound {
		logger.WithFields(log.Fields{
			"contentType": contentType,
		}).Fatal("Unrecognised content type")
	}
	return media
}

func GetPayloadPrefix(media http.ContentType) []byte {
	payloadPrefix, ok := payloadPrefixes[media]

	if !ok {
		logger.WithFields(log.Fields{
			"contentType": media,
		}).Fatal("No payload prefix found")
	}

	return payloadPrefix
}
