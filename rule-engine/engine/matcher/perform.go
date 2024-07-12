// Under the Apache-2.0 License
package matcher

import (
	"github.com/groboclown/qazaar-testing/rule-engine/engine/obj"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/descriptor"
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

	// NOT operations should operate on just one item, but for optimization,
	// this will just run the matcher with an AND rule, then return the
	// opposite of that.
	if operation == srule.NotCollection {
		return !IsMatch(obj, srule.AndCollection, matcher)
	}

	// AND and OR logic is identical except for the matching condition
	// for early loop exit.
	earlyExitCondition := false
	if operation == srule.OrCollection {
		earlyExitCondition = true
	}

	for _, c := range matcher.Collection {
		if earlyExitCondition == IsMatch(obj, c.Operation, matcher) {
			return earlyExitCondition
		}
	}

	for _, c := range matcher.Contains {
		if earlyExitCondition == IsContainsMatch(obj, &c) {
			return earlyExitCondition
		}
	}
	return !earlyExitCondition
}

func IsContainsMatch(
	obj *obj.EngineObj,
	contains *srule.ContainsMatcher,
) bool {
	if contains == nil || obj == nil {
		return false
	}

	val, distinct := obj.Value(contains.Key)
	if contains.Distinct && !distinct {
		// Only perform the extra make-it-distinct logic if the key isn't already distinct.
		val = val.Distinct()
	}
	if contains.Count {
		val = val.CountValue()
	}

	if listMatch(contains.Operation, val.Number, contains.Checks.Numeric, numericCheckConst) {
		return true
	}
	if listMatch(contains.Operation, val.Text, contains.Checks.Text, stringCheckConst) {
		return true
	}

	return false
}

func listMatch[T descriptor.DescriptorValueTypes, C srule.NumericBoundsCheck | srule.StringCheck](
	operation srule.ContainsOperation,
	values []T,
	cmp []C,
	check vcheck[T, C],
) bool {
	cmpCount := len(cmp)
	valueCount := len(values)

	// There are some special edge cases that we can check for and exit early.

	if cmpCount == 0 || valueCount == 0 {
		// Most likely an empty value or incorrect type check.
		return false
	}

	if operation == srule.ContainsExactly && valueCount != cmpCount {
		// This special case is worth investigating.
		// For 'exactly', it must be a 1-to-1 match between a comparison and a value.
		// So if the value count != comparison count, then this cannot happen.
		// Don't even need to check the values.
		return false
	}

	if operation == srule.ContainsSome || (operation == srule.ContainsAll && cmpCount == 1) {
		// Exit on first match.
		for _, c := range cmp {
			for _, v := range values {
				if check.check(v, &c) {
					return true
				}
			}
		}
		return false
	}

	// In all other cases, a comparison count of checks vs. matched vs. values
	// handles the other operations.

	// Perform the matching over the values.
	// Each value can be matched at most 1 time.
	matchIdx := make(map[int]bool)
	matchCount := 0
	for _, c := range cmp {
		for i, v := range values {
			if !matchIdx[i] && check.check(v, &c) {
				matchIdx[i] = true
				matchCount++
			}
		}
	}

	switch operation {
	case srule.ContainsAll:
		// All the comparisons must match against a value.
		return len(matchIdx) == cmpCount
	case srule.ContainsExactly:
		// The match count, comparison count, and value count must all equal.
		// However, the above comparison already ensures that value count == comparison count.
		return matchCount == cmpCount
	case srule.ContainsOnly:
		// Every value must have been matched.
		return matchCount == valueCount
	}
	// Fail with a runtime error?
	return false
}

type vcheck[T descriptor.DescriptorValueTypes, C srule.NumericBoundsCheck | srule.StringCheck] interface {
	check(v T, c *C) bool
}

type numericCheck int

func (n numericCheck) check(v float64, c *srule.NumericBoundsCheck) bool {
	return v >= c.Min && v <= c.Max
}

var numericCheckConst = numericCheck(0)

type stringCheck int

func (n stringCheck) check(v string, c *srule.StringCheck) bool {
	return c.Matches(v)
}

var stringCheckConst = stringCheck(0)
