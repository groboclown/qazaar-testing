// Under the Apache-2.0 License
package sog

import (
	"context"
	"strconv"
	"strings"

	"github.com/groboclown/qazaar-testing/rule-engine/engine/matcher"
	"github.com/groboclown/qazaar-testing/rule-engine/engine/obj"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/srule"
)

// SogBuilder builds SOG objects for a single group definition.
//
// By its nature, it must scan all objects.
// Because SOG definitions create new objects, this incrementally builds with
// new objects.
// In order to prevent a recursion issue, it inserts a parent marker in the
// constructed objects.
type SogBuilder struct {
	id      string
	rule    *srule.Group
	byId    map[string]*sogInstanceBuilder
	factory obj.ObjFactory
}

// NewBuilder creates a new SogBuilder instance.
func NewBuilder(rule *srule.Group, factory obj.ObjFactory) *SogBuilder {
	if rule == nil {
		return nil
	}
	return &SogBuilder{
		id:      rule.Id,
		rule:    rule,
		byId:    make(map[string]*sogInstanceBuilder),
		factory: factory,
	}
}

// Seal closes off all current builder synthetic instances from additional members, and returns them.
//
// This returns the engine object representation of the SOG instances.
// This can be safely called multiple times.
func (s *SogBuilder) Seal() []SogInstance {
	ret := make([]SogInstance, len(s.byId))
	i := 0
	for _, v := range s.byId {
		ret[i] = &sogInstance{
			members: v.members,
			obj:     v.seal(s.factory, s.id, s.rule.Alterations),
			group:   s.rule,
		}
		i++
	}
	return ret
}

// Reset clears out the current list of known SOG instances.
func (s *SogBuilder) Reset() {
	s.byId = make(map[string]*sogInstanceBuilder)
}

type SogBuilderAddResult int

const (
	NoMatch SogBuilderAddResult = iota
	Created
	Added
	Recursion
)

// Add adds the object to the SOG, and returns the instance + if it's a newly created instance.
//
// If the object does not match this builder, then it returns NoMatch.
func (s *SogBuilder) Add(o *obj.EngineObj) SogBuilderAddResult {
	if o == nil {
		return NoMatch
	}
	if !matcher.IsMatch(o, s.rule.Matchers) {
		return NoMatch
	}
	shared := groupSharedValues(s.rule, o)
	si := s.matching(shared)
	ret := Added
	if si == nil {
		ret = Created

		// Build a unique identifier
		baseId := matchGroupId(shared)
		id := baseId
		var idx int64 = 0
		for {
			if _, ok := s.byId[id]; !ok {
				break
			}
			idx++
			id = baseId + ":" + strconv.FormatInt(idx, 10)
		}

		// Add the new instance to the internal record.
		si = newSogInstance(id, shared)
		s.byId[si.id] = si
	} else {
		if si.IsRecursion(o, context.Background()) {
			// Do not add the item.
			return Recursion
		}
	}
	si.AddMember(o)
	return ret
}

// matchGroup finds the sog instance associated with the engine object.
//
// If the engine object does not match this rule, then it returns nil.
// If this object creates a new instance, then it returns that instance.
func (s *SogBuilder) matching(shared map[string]obj.DescriptorValues) *sogInstanceBuilder {
	// No identifier derived from the shared keys + values is guaranteed unique, so
	// we must look through item by item.
	// Alternatively, we could devise a quick lookup scheme by values, because we
	// know the keys are the same.  That's a future optimization, though.
	for _, si := range s.byId {
		if si.matches(shared) {
			return si
		}
	}
	return nil
}

func groupSharedValues(rule *srule.Group, o *obj.EngineObj) map[string]obj.DescriptorValues {
	ret := make(map[string]obj.DescriptorValues)
	for _, k := range rule.KeySharedValues {
		val, _ := o.Value(k)
		ret[k] = val
	}
	return ret
}

// matchGroupId creates an identifier for the group instance based on the shared values.
func matchGroupId(shared map[string]obj.DescriptorValues) string {
	parts := make([]string, 0, len(shared)*4)

	for key, value := range shared {
		parts = append(parts, "&", key, ":")
		joiner := ""
		for _, n := range value.Number {
			parts = append(parts, joiner, strconv.FormatFloat(n, 'f', -1, 64))
			joiner = "|"
		}
	}

	return strings.Join(parts, "")
}
