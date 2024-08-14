// Under the Apache-2.0 License
package obj

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/descriptor"
)

// Value returns a descriptor value list for the given key.
//
// All objects contain every key, though the default value is an empty list.
func (o *EngineObj) Value(key string) (DescriptorValues, bool) {
	// This takes up potentially lots of extra memory for fun memory swapping times.
	// Could potentially use a memory cache to reuse the values, but before that
	// extra effort, it needs some solid profiling to see if caching would pay off.
	if e, ok := o.Enum[key]; ok {
		return DescriptorValues{
			Text:   e.List(),
			Number: nil,
		}, e.IsDistinct()
	}
	if f, ok := o.Free[key]; ok {
		return DescriptorValues{
			Text:   f.List(),
			Number: nil,
		}, f.IsDistinct()
	}
	if n, ok := o.Numeric[key]; ok {
		return DescriptorValues{
			Text:   nil,
			Number: n.List(),
		}, n.IsDistinct()
	}
	return DescriptorValues{}, false
}

// Count returns the number of values for the key in this object.
//
// All objects contain every key, though the default value is an empty list.
func (o *EngineObj) Count(key string) int {
	if e, ok := o.Enum[key]; ok {
		return e.Count()
	}
	if f, ok := o.Free[key]; ok {
		return f.Count()
	}
	if n, ok := o.Numeric[key]; ok {
		return n.Count()
	}
	return 0
}

// CountValue turns the key into a value object containing the count as the only value (numeric).
func (o *EngineObj) CountValue(key string) DescriptorValues {
	return DescriptorValues{Number: []float64{float64(o.Count(key))}}
}

func (o *EngineObj) String() string {
	attribs := make([]string, 0)
	for k, e := range o.Enum {
		attribs = append(attribs, fmt.Sprintf("%s:%s", k, strings.Join(e.List(), ",")))
	}
	for k, f := range o.Free {
		attribs = append(attribs, fmt.Sprintf("%s:%s", k, strings.Join(f.List(), ",")))
	}
	for k, n := range o.Numeric {
		vals := n.List()
		vs := make([]string, len(vals))
		for i, v := range vals {
			vs[i] = strconv.FormatFloat(v, 'f', 4, 64)
		}
		attribs = append(attribs, fmt.Sprintf("%s:%s", k, strings.Join(vs, ",")))
	}
	return o.Id + "(" + o.Source.String() + ")" + "{" + strings.Join(attribs, ";") + "}"
}

// Distinct turns the values into a distinct list of values.
func (d DescriptorValues) Distinct() DescriptorValues {
	var num []float64 = nil
	if d.Number != nil {
		num = descriptor.DistinctValueArray(d.Number, descriptor.FloatTransform)
	}
	var txt []string = nil
	if d.Text != nil {
		txt = descriptor.DistinctValueArray(d.Text, descriptor.StringTransform)
	}
	return DescriptorValues{Text: txt, Number: num}
}

// CountValue turns the key into a value object containing the count as the only value (numeric).
func (d DescriptorValues) CountValue() DescriptorValues {
	return DescriptorValues{Number: []float64{float64(len(d.Number) + len(d.Text))}}
}

func (d DescriptorValues) Count() int {
	return len(d.Number) + len(d.Text)
}

func (s ObjSource) String() string {
	ret := ""
	if s.Construct != nil {
		ret = fmt.Sprintf("[%s] ", *s.Construct)
	}
	if len(s.Parents) > 0 {
		first := true
		ret = ret + "("
		for _, p := range s.Parents {
			if first {
				first = false
			} else {
				ret = ret + ", "
			}
			ret = ret + p.String()
		}
		return ret + ")"
	}

	if len(s.Source) <= 0 {
		return ret + "?"
	}
	ver := s.Source[0].Ver()
	a := s.Source[0].A()
	rep := s.Source[0].Rep()

	ret = ret + s.Source[0].Loc()
	if rep != "" {
		ret = rep + ":" + ret
	}
	if a != nil && *a != "" {
		ret = rep + "#" + *a
	}
	if ver != nil && *ver != "" {
		ret = rep + "@" + *ver
	}
	return ret
}
