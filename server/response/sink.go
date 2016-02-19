package response

import (
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
)

// Sink is a part of the response which is returned by sinks' action.
type Sink struct {
	NodeType string      `json:"node_type"`
	Name     string      `json:"name"`
	State    string      `json:"state"`
	Status   data.Map    `json:"status,omitempty"`
	Meta     interface{} `json:"meta,omitempty"`
}

// NewSink returns the result of the sink node. It generates status and
// meta information if detailed argument is true.
func NewSink(sn core.SinkNode, detailed bool) *Sink {
	s := &Sink{
		NodeType: core.NTSink.String(),
		Name:     sn.Name(),
		State:    sn.State().Get().String(),
	}

	if detailed {
		s.Status = sn.Status()
		s.Meta = sn.Meta()
	}
	return s
}
