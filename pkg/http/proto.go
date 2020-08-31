package http

// Protocol represents an HTTP protocol
type Protocol string

// Valid HTTP protocols
const (
	HTTP10 Protocol = "HTTP/1.0"
	HTTP11 Protocol = "HTTP/1.1"
	HTTP20 Protocol = "HTTP/2.0"
)

func (p Protocol) String() string {
	return string(p)
}
