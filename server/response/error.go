package response

import (
	"gopkg.in/sensorbee/sensorbee.v0/data"
)

// Error has error response information of a request.
type Error struct {
	Code      string   `json:"code"`
	Message   string   `json:"message"`
	RequestID string   `json:"request_id"`
	Meta      data.Map `json:"meta"`
}
