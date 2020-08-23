package http

type HttpVerb int

const (
	Head HttpVerb = iota
	Get
	Post
	Delete
	Put
	Connect
	Options
	Trace
	Patch
)

func (v HttpVerb) String() string {
	switch v {
	case Head:
		return "HEAD"
	case Get:
		return "GET"
	case Post:
		return "POST"
	case Delete:
		return "DELETE"
	case Put:
		return "PUT"
	case Connect:
		return "CONNECT"
	case Options:
		return "OPTIONS"
	case Trace:
		return "TRACE"
	case Patch:
		return "PATCH"
	}

	return "UNRECOGNISED"
}

func (v HttpVerb) hasBody() bool {
	switch v {
	case Put, Post, Patch:
		return true
	}

	return false
}
