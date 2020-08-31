package http

// verb represents and HTTP verb/action
type verb string

// Valid HTTP verbs
const (
	Head    verb = "HEAD"
	Get     verb = "GET"
	Post    verb = "POST"
	Delete  verb = "DELETE"
	Put     verb = "PUT"
	Connect verb = "CONNECT"
	Options verb = "OPTIONS"
	Trace   verb = "TRACE"
	Patch   verb = "PATCH"
)

func (v verb) String() string {
	return string(v)
}

func (v verb) hasBody() bool {
	switch v {
	case Put, Post, Patch:
		return true
	}

	return false
}
