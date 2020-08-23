package http

const (
	ApplicationJson ContentType = "application/json"
)

type ContentType string

func (c ContentType) String() string {
	return string(c)
}
