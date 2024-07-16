// Under the Apache-2.0 License
package sog

import (
	"context"

	"github.com/groboclown/qazaar-testing/rule-engine/engine/concurrent"
	"github.com/groboclown/qazaar-testing/rule-engine/engine/obj"
)

// IsRecursion checks if the engine object is somewhere in the parent list.
func (s *SogInstance) IsRecursion(o *obj.EngineObj, ctx context.Context) bool {
	if s == nil || o == nil {
		return false
	}
	initial := make([]*obj.ObjSource, len(s.members))
	for i, m := range s.members {
		if m == o {
			return true
		}
		initial[i] = &m.Source
	}
	ret := <-concurrent.RunEarlyExit(recCheck{&o.Source}, ctx, initial)
	return ret == &recYes
}

type recCheck struct {
	against *obj.ObjSource
}

var recYes = true

func (r recCheck) Perform(v *obj.ObjSource, otherWorkers chan<- *obj.ObjSource) (*bool, bool) {
	if r.against == v {
		return &recYes, true
	}
	for _, p := range v.Parents {
		otherWorkers <- p
	}
	return nil, false
}
