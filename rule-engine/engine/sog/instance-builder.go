// Under the Apache-2.0 License
package sog

import (
	"github.com/groboclown/qazaar-testing/rule-engine/engine/obj"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/descriptor"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/srule"
)

type sogInstanceBuilder struct {
	id       string
	shared   map[string]obj.DescriptorValues
	members  []*obj.EngineObj
	instance *obj.EngineObj
	sealed   bool
}

func newSogInstance(
	id string,
	shared map[string]obj.DescriptorValues,
) *sogInstanceBuilder {
	return &sogInstanceBuilder{
		id:       id,
		shared:   shared,
		members:  make([]*obj.EngineObj, 0),
		instance: nil,
		sealed:   false,
	}
}

// AddMember adds a new engine object to the member list.
//
// If the instance was previously sealed, this panics.
func (s *sogInstanceBuilder) AddMember(o *obj.EngineObj) {
	if s == nil || o == nil {
		return
	}
	if s.sealed {
		panic("attempted to add a member to a sealed instance")
	}
	s.instance = nil
	s.members = append(s.members, o)
}

// seal closes off the SOG instance, creating the synthetic object.
//
// This prevents adding new members to the instance.  This can safely
// be called multiple times.
func (s *sogInstanceBuilder) seal(
	factory obj.ObjFactory,
	groupId string,
	alterations []srule.Alteration,
) *obj.EngineObj {
	if s == nil {
		return nil
	}
	s.sealed = true

	if s.instance == nil {
		s.instance = asFinalizedObj(
			factory.FromGroup(s.members, groupId),
			s.shared,
			alterations,
		)
	}
	return s.instance
}

func (s *sogInstanceBuilder) Sealed() bool {
	return s.sealed
}

// matches checks if the matched keys passed match the shared values stored.
//
// This assumes that the shared values have the same keys as the sog instance.
func (s *sogInstanceBuilder) matches(shared map[string]obj.DescriptorValues) bool {
	for k, da := range shared {
		db := s.shared[k]
		if !matchesDescriptorValues(&da, &db) {
			return false
		}
	}
	return true
}

func matchesDescriptorValues(a *obj.DescriptorValues, b *obj.DescriptorValues) bool {
	// Early exit conditions.
	if a == b {
		return true
	}
	if len(a.Number) != len(b.Number) || len(a.Text) != len(b.Text) {
		return false
	}
	if !matchesTypedValues(a.Number, b.Number) {
		return false
	}
	return matchesTypedValues(a.Text, b.Text)
}

func matchesTypedValues[T descriptor.DescriptorValueTypes](a []T, b []T) bool {
	// Note: order doesn't matter.
	// But, for performance reasons, this could be faster for distinct values.
	missingA := make(map[T]int)
	for _, v := range a {
		i, ok := missingA[v]
		if !ok {
			i = 0
		}
		missingA[v] = i + 1
	}
	for _, v := range b {
		c, ok := missingA[v]
		if !ok || c <= 0 {
			return false
		}
		missingA[v] = c - 1
	}
	for _, c := range missingA {
		if c > 0 {
			return false
		}
	}
	return true
}
