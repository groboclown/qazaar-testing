// Under the Apache-2.0 License
package sog

import (
	"github.com/groboclown/qazaar-testing/rule-engine/engine/obj"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/srule"
)

type SogInstance interface {
	Members() []*obj.EngineObj
	Obj() *obj.EngineObj
	Group() *srule.Group
}

type sogInstance struct {
	members []*obj.EngineObj
	obj     *obj.EngineObj
	group   *srule.Group
}

func (s *sogInstance) Members() []*obj.EngineObj {
	return s.members
}

func (s *sogInstance) Obj() *obj.EngineObj {
	return s.obj
}

func (s *sogInstance) Group() *srule.Group {
	return s.group
}
