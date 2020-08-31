package http

// Valid HTTP content types
const (
	ApplicationJSON               ContentType = "application/json"
	ApplicationXML                ContentType = "application/xml"
	ApplicationXWWWFormURLEncoded ContentType = "application/x-www-form-urlencoded"
	ContentTypeNotFound           ContentType = "NOT_FOUND"
)

var (
	contentTypes = func() map[string]ContentType {
		types := []ContentType{
			ApplicationJSON,
			ApplicationXML,
			ApplicationXWWWFormURLEncoded,
		}

		mapped := make(map[string]ContentType)

		for _, contentType := range types {
			mapped[contentType.String()] = contentType
		}

		return mapped
	}()
)

// ContentType represents an HTTP content type
type ContentType string

func (c ContentType) String() string {
	return string(c)
}

// ToContentType converts a string into a representative ContentType. In the event no ContentType
// corresponding to the string the null object ContentTypeNotFound is returned. In this case a
// ContentType may be created manually if desired
func ToContentType(contentType string) ContentType {
	val, ok := contentTypes[contentType]

	if ok {
		return val
	}

	return ContentTypeNotFound
}
