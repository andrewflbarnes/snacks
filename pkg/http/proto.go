package http

type HttpProto int

const (
	Http10 HttpProto = iota
	Http11
	Http20
)

func (p HttpProto) String() string {
	switch p {
	case Http10:
		return "HTTP/1.0"
	case Http11:
		return "HTTP/1.1"
	case Http20:
		return "HTTP/2.0"
	}

	return "UNRECOGNISED"
}
