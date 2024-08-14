// Match convergence rules.
//
// Under the Apache-2.0 License
package runner

import (
	"github.com/groboclown/qazaar-testing/rule-engine/engine/obj"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/descriptor"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/sont"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/srule"
)

func MatchConvergence(
	groupId string,
	objs []*obj.EngineObj,
	c *srule.Convergence,
	ont *sont.AllowedDescriptors,
) *ConvProblem {
	d := ont.Find(c.Key)
	if d == nil || len(objs) <= 0 {
		// FIXME No such key; should be a Problem instead.
		return nil
	}
	groups := findGroup(objs, d, c.Distinct)

	switch c.Requires {
	case srule.AllMatch:
		return convMatchRes(groupId, c, matchAll(groups))
	case srule.Disjoint:
		return convMatchRes(groupId, c, matchDisjoint(objs, groups))
	default:
		// FIXME Should be a Problem
		return nil
	}
}

func matchAll(ml []MemberValues) []MemberValues {
	if len(ml) == 1 {
		// 1 value, which means that they all share the same value.
		return nil
	}
	return ml
}

func matchDisjoint(objs []*obj.EngineObj, ml []MemberValues) []MemberValues {
	if len(objs) == len(ml) {
		// Each object is its own group, so they are all distinct from each other.
		return nil
	}
	return ml
}

func newMemberValues(o *obj.EngineObj, v obj.DescriptorValues) *MemberValues {
	return &MemberValues{
		Members: []*obj.EngineObj{o},
		Text:    v.Text,
		Number:  v.Number,
	}
}

// dMemberValues allows for faster group checking with an occurrence count map.
type dMemberValues struct {
	m         *MemberValues
	tm        map[string]int
	nm        map[float64]int
	distinct  bool
	textTform descriptor.DescriptorValueTypeTransform[string]
}

func newDMemberValues(
	o *obj.EngineObj,
	v obj.DescriptorValues,
	d *sont.TypedDescriptor,
	forceDistinct bool,
) *dMemberValues {
	ret := &dMemberValues{
		m:         newMemberValues(o, v),
		tm:        make(map[string]int),
		nm:        make(map[float64]int),
		distinct:  forceDistinct || d.IsDistinct(),
		textTform: descriptor.StringLowerTransform,
	}
	if d.IsCaseSensitive() {
		ret.textTform = descriptor.StringTransform
	}
	loadMaps(ret.tm, ret.nm, v, ret.textTform)
	return ret
}

func loadMaps(
	tm map[string]int,
	nm map[float64]int,
	v obj.DescriptorValues,
	tf descriptor.DescriptorValueTypeTransform[string],
) {
	for _, t := range v.Text {
		// This should work.  The default getter should return 0.
		tm[tf.Transform(t)]++
	}
	for _, n := range v.Number {
		nm[n]++
	}
}

func (d *dMemberValues) matchAndAdd(o *obj.EngineObj, v obj.DescriptorValues) bool {
	if d.match(v) {
		d.m.Members = append(d.m.Members, o)
		return true
	}
	return false
}

func (d *dMemberValues) match(v obj.DescriptorValues) bool {
	if len(v.Number) != len(d.nm) || len(v.Text) != len(d.tm) {
		return false
	}
	if d.distinct {
		for _, n := range v.Number {
			if d.nm[n] < 1 {
				return false
			}
		}
		for _, t := range v.Text {
			if d.tm[d.textTform.Transform(t)] < 1 {
				return false
			}
		}
		return true
	}

	vt := make(map[string]int)
	vn := make(map[float64]int)
	loadMaps(vt, vn, v, d.textTform)
	for k, c := range vt {
		if d.tm[k] != c {
			return false
		}
	}
	for k, c := range vn {
		if d.nm[k] != c {
			return false
		}
	}
	return true
}

// findGroup turns the objects into groups of shared values, based on the type descriptor.
func findGroup(objs []*obj.EngineObj, d *sont.TypedDescriptor, distinct bool) []MemberValues {
	key := d.KeyName()
	dGroups := make([]*dMemberValues, 0)
	groups := make([]MemberValues, 0)
	for _, o := range objs {
		v, _ := o.Value(key)
		missing := true
		for _, g := range dGroups {
			if g.matchAndAdd(o, v) {
				missing = false
				break
			}
		}
		if missing {
			dmv := newDMemberValues(o, v, d, distinct)
			dGroups = append(dGroups, dmv)
			// While this makes a copy of the object, it should be
			// a shallow copy, meaning the list of members is shared,
			// and that's the important part of this logic.
			groups = append(groups, *dmv.m)
		}
	}
	return groups
}

func convMatchRes(
	groupId string,
	conv *srule.Convergence,
	res []MemberValues,
) *ConvProblem {
	if res == nil {
		return nil
	}
	return &ConvProblem{
		GroupId:    groupId,
		Mismatched: res,
		Conv:       conv,
	}
}
