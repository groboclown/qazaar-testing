// Matcher violation data model.
//
// The string generation here has very English-y assumptions.  It will need
// better templatization for future localization efforts.
//
// Under the Apache-2.0 License
package matcher

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/groboclown/qazaar-testing/rule-engine/engine/obj"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/srule"
)

type MatcherMismatch struct {
	Obj        *obj.EngineObj
	Collection *MismatchCollection
	Contains   *MismatchContains
}

func (m MatcherMismatch) String() string {
	ret := fmt.Sprintf("Mismatch for %s: ", m.Obj.String())
	if m.Collection != nil {
		ret += m.Collection.String(m.Obj)
	}
	if m.Contains != nil {
		ret += m.Contains.String(m.Obj)
	}
	return ret
}

type MismatchCollection struct {
	Operation srule.CollectionOperation
	Matcher   *srule.MatchingDescriptorSet
}

func (m MismatchCollection) String(parent *obj.EngineObj) string {
	join := ""
	switch m.Operation {
	case srule.AndCollection:
		join = " AND "
	case srule.OrCollection:
		join = " OR "
	case srule.NotCollection:
		return "NOT " + (MismatchCollection{
			Operation: srule.AndCollection,
			Matcher:   m.Matcher,
		}).String(parent)
	default:
		join = " <UNSUPPORTED> "
	}
	parts := make([]string, 0)
	for _, p := range m.Matcher.Collection {
		parts = append(parts, (MismatchCollection{
			Operation: p.Operation,
			Matcher:   p.Matchers,
		}).String(parent))
	}
	for _, p := range m.Matcher.Contains {
		parts = append(parts, (MismatchContains{Contains: &p}).String(parent))
	}
	return "(" + strings.Join(parts, join) + ")"
}

type MismatchContains struct {
	Contains *srule.ContainsMatcher
}

func (m MismatchContains) String(parent *obj.EngineObj) string {
	ret := m.Contains.Key
	if m.Contains.Distinct {
		ret = "distinct " + ret
	}
	if m.Contains.Count {
		ret = "count of " + ret
	}

	switch m.Contains.Operation {
	case srule.ContainsAll:
		ret += " contains all "
	case srule.ContainsExactly:
		ret += " contains exactly "
	case srule.ContainsOnly:
		ret += " contains only "
	case srule.ContainsSome:
		ret += " contains some "
	default:
		ret += " <UNSUPPORTED> "
	}

	first := true
	for _, c := range m.Contains.Checks.Numeric {
		if first {
			first = false
		} else {
			ret += ", "
		}
		ret += fmt.Sprintf("range [%f, %f]", c.Min, c.Max)
	}
	for _, c := range m.Contains.Checks.Text {
		if first {
			first = false
		} else {
			ret += ", "
		}
		ret += "'" + c.R.String() + "'"
	}
	ret += " but has ("
	val, _ := parent.Value(m.Contains.Key)
	first = true
	for _, v := range val.Number {
		if first {
			first = false
		} else {
			ret += ", "
		}
		ret += strconv.FormatFloat(v, 'f', 4, 64)
	}
	for _, v := range val.Text {
		if first {
			first = false
		} else {
			ret += ", "
		}
		ret += "'" + v + "'"
	}
	return ret + ")"
}
