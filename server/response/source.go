package response

import (
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
)

// Source is a part of the response which is returned by sources' action.
type Source struct {
	NodeType string      `json:"node_type"`
	Name     string      `json:"name"`
	State    string      `json:"state"`
	Status   data.Map    `json:"status,omitempty"`
	Meta     interface{} `json:"meta,omitempty"`
}

// NewSource returns the result of the source node. It generates status and
// meta information if detailed argument is true.
func NewSource(sn core.SourceNode, detailed bool) *Source {
	s := &Source{
		NodeType: core.NTSource.String(),
		Name:     sn.Name(),
		State:    sn.State().Get().String(),
	}

	if detailed {
		s.Status = sn.Status()
		s.Meta = sn.Meta()
	}
	return s
}
