// Under the Apache-2.0 License
package matcher

import (
	"github.com/groboclown/qazaar-testing/rule-engine/engine/obj"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/srule"
)

func IsMatch(
	obj *obj.EngineObj,
	operation srule.CollectionOperation,
	matcher *srule.MatchingDescriptorSet,
) bool {
	if matcher == nil || obj == nil {
		return false
	}
	for _, c := range matcher.Collection {
		if !IsMatch(obj, c.Operation, matcher) {
			return false
		}
	}
	for _, c := range matcher.Contains {
		if !IsContainsMatch(obj, &c) {
			return false
		}
	}
	return true
}

func IsContainsMatch(
	obj *obj.EngineObj,
	contains *srule.ContainsMatcher,
) bool {
	if contains == nil || obj == nil {
		return false
	}
	if contains.Count {
		panic("not implemented")
	}
	if contains.Distinct {
		panic("not implemented")
	}
	panic("not implemented")
}
