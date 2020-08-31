package http

const (
	ApplicationJson               ContentType = "application/json"
	ApplicationXWWWFormUrlEncoded ContentType = "application/x-www-form-urlencoded"
)

type ContentType string

func (c ContentType) String() string {
	return string(c)
}
