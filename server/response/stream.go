package response

import (
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// Stream is a part of the response which is returned by streams' action.
type Stream struct {
	NodeType string      `json:"node_type"`
	Name     string      `json:"name"`
	State    string      `json:"state"`
	Status   data.Map    `json:"status,omitempty"`
	Meta     interface{} `json:"meta,omitempty"`
}

// NewStream returns the result of the box node. It generates status and
// meta information if detailed argument is true.
func NewStream(bn core.BoxNode, detailed bool) *Stream {
	s := &Stream{
		NodeType: core.NTBox.String(),
		Name:     bn.Name(),
		State:    bn.State().Get().String(),
	}

	if detailed {
		s.Status = bn.Status()
		s.Meta = bn.Meta()
	}
	return s
}
