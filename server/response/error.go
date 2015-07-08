package response

import (
	"pfi/sensorbee/sensorbee/data"
)

// Error has error response information of a request.
type Error struct {
	Code      string   `json:"code"`
	Message   string   `json:"message"`
	RequestID string   `json:"request_id"`
	Meta      data.Map `json:"meta"`
}
