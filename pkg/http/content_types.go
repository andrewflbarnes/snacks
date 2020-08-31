package http

const (
	ApplicationJSON               ContentType = "application/json"
	ApplicationXWWWFormURLEncoded ContentType = "application/x-www-form-urlencoded"
	ContentTypeNotFound           ContentType = "NOT_FOUND"
)

var (
	contentTypes = func() map[string]ContentType {
		types := []ContentType{
			ApplicationJSON,
			ApplicationXWWWFormURLEncoded,
		}

		mapped := make(map[string]ContentType)

		for _, contentType := range types {
			mapped[contentType.String()] = contentType
		}

		return mapped
	}()
)

type ContentType string

func (c ContentType) String() string {
	return string(c)
}

func ToContentType(contentType string) ContentType {
	val, ok := contentTypes[contentType]

	if ok {
		return val
	} else {
		return ContentTypeNotFound
	}
}
