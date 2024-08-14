// Under the Apache-2.0 License
package sog

import (
	"context"

	"github.com/groboclown/qazaar-testing/rule-engine/engine/concurrent"
	"github.com/groboclown/qazaar-testing/rule-engine/engine/obj"
)

// IsRecursion checks if the engine object is somewhere in the parent list.
func (s *sogInstanceBuilder) IsRecursion(o *obj.EngineObj, ctx context.Context) bool {
	if s == nil || o == nil {
		return false
	}
	q, r := concurrent.NewEarlyExit[*obj.ObjSource](0)
	go fetchAllSources(q, s.members)
	for {
		select {
		case v, ok := <-r.Reader():
			if &o.Source == v {
				r.Stop()
				return true
			}
			if !ok {
				r.Stop()
				return false
			}
		case <-ctx.Done():
			r.Stop()
			return false
		}
	}
}

func fetchAllSources(q concurrent.EarlyExitQueue[*obj.ObjSource], members []*obj.EngineObj) {
	stack := make([]*obj.ObjSource, len(members))
	for i, m := range members {
		stack[i] = &m.Source
	}

	for len(stack) > 0 {
		curr := stack[0]
		stack = stack[1:]
		if _, stopped := q.QueueAll(curr.Parents...); stopped {
			return
		}
		stack = append(stack, curr.Parents...)
	}
	q.Finished()
}
