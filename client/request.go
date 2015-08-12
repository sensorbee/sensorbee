package client

// Method has a type of HTTP methods.
type Method int

const (
	// Get represents GET method.
	Get Method = iota

	// Post represents POST method.
	Post

	// Put represents PUT method.
	Put

	// Delete represents DELETE method.
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
