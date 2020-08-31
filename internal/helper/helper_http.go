package helper

import (
	"github.com/andrewflbarnes/snacks/pkg/http"
)

var (
	ApplicationJsonPrefix               = newMediaPrefixImpl(http.ApplicationJson, []byte(`{"a":"`))
	ApplicationXWWWFormUrlencodedPrefix = newMediaPrefixImpl(http.ApplicationXWWWFormUrlEncoded, []byte(`a=`))
)

type MediaPrefix interface {
	Name() string
	ContentType() http.ContentType
	Prefix() []byte
}

type mediaPrefixImpl struct {
	contentType http.ContentType
	bodyPrefix  []byte
}

func (m mediaPrefixImpl) Name() string {
	return m.contentType.String()
}

func (m mediaPrefixImpl) ContentType() http.ContentType {
	return m.contentType
}

func (m mediaPrefixImpl) Prefix() []byte {
	return m.bodyPrefix
}

func newMediaPrefixImpl(contentType http.ContentType, bodyPrefix []byte) MediaPrefix {
	return mediaPrefixImpl{
		contentType: contentType,
		bodyPrefix:  bodyPrefix,
	}
}
