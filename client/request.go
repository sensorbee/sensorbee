package client

type Method int

const (
	Get Method = iota
	Post
	Put
	Delete
)

func (r Method) String() string {
	switch r {
	case Get:
		return "GET"
	case Post:
		return "POST"
	case Put:
		return "PUT"
	case Delete:
		return "DELETE"
	default:
		return "unknown"
	}
}
