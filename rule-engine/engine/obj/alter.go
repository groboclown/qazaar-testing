// Under the Apache-2.0 License
package obj

import (
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/descriptor"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/sources"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/sont"
)

type engineObjBuilder struct {
	source ObjSource
	ont    *sont.AllowedDescriptors

	numeric map[string]descriptor.DescriptorValueBuilder[float64]
	enum    map[string]descriptor.DescriptorValueBuilder[string]
	free    map[string]descriptor.DescriptorValueBuilder[string]
}

func newObjBuilder(
	parents []*ObjSource,
	construct *string,
	source []sources.Source,
	ont *sont.AllowedDescriptors,
) *engineObjBuilder {
	return &engineObjBuilder{
		source: ObjSource{
			Parents:   parents,
			Construct: construct,
			Source:    source,
		},
		ont:     ont,
		numeric: make(map[string]descriptor.DescriptorValueBuilder[float64]),
		enum:    make(map[string]descriptor.DescriptorValueBuilder[string]),
		free:    make(map[string]descriptor.DescriptorValueBuilder[string]),
	}
}

func (o *EngineObj) Alter() EngineObjBuilder {
	if o == nil {
		return nil
	}

	return &engineObjBuilder{
		source:  o.Source,
		ont:     o.ont,
		numeric: mutateDescriptorMap(o.Numeric),
		enum:    mutateDescriptorMap(o.Enum),
		free:    mutateDescriptorMap(o.Free),
	}
}

// Add adds new values to the key's value.
func (o *engineObjBuilder) Add(key string, val DescriptorValues) {
	if o == nil {
		return
	}
	alt := &alterAction{val}
	o.alter(key, false, alt.addFloat, alt.addString)
}

// Add adds new values to the key's value, unless that value already exists in the key.
func (o *engineObjBuilder) AddDistinct(key string, val DescriptorValues) {
	if o == nil {
		return
	}
	alt := &alterAction{val}
	o.alter(key, false, alt.addDistinctFloat, alt.addDistinctString)
}

func (o *engineObjBuilder) Remove(key string, val DescriptorValues) {
	if o == nil {
		return
	}
	alt := &alterAction{val}
	o.alter(key, true, alt.removeFloat, alt.removeString)
}

func (o *engineObjBuilder) RemoveDistinct(key string, val DescriptorValues) {
	if o == nil {
		return
	}
	alt := &alterAction{val}
	o.alter(key, true, alt.removeDistinctFloat, alt.removeDistinctString)
}

func (o *engineObjBuilder) Set(key string, val DescriptorValues) {
	if o == nil {
		return
	}
	alt := &alterAction{val}
	o.alter(key, true, alt.setFloat, alt.setString)
}

type alterFloatFunc func(descriptor.DescriptorValueBuilder[float64])
type alterStringFunc func(descriptor.DescriptorValueBuilder[string])

func (o *engineObjBuilder) alter(key string, removes bool, fn alterFloatFunc, fs alterStringFunc) {
	if n, ok := o.numeric[key]; ok {
		fn(n)
		return
	}
	if e, ok := o.enum[key]; ok {
		fs(e)
		return
	}
	if f, ok := o.free[key]; ok {
		fs(f)
		return
	}

	if removes {
		// Removal should not add an attribute, as that will end up just taking time
		// and bloating memory.
		return
	}

	if nd := o.ont.Numeric(key); nd != nil {
		n := descriptor.NewNumericBuilder(nd.IsDistinct())
		fn(n)
		o.numeric[key] = n
		return
	}
	if ed := o.ont.Enum(key); ed != nil {
		e := descriptor.NewTextBuilder(ed.IsDistinct(), true)
		fs(e)
		o.enum[key] = e
		return
	}
	if fd := o.ont.Free(key); fd != nil {
		f := descriptor.NewTextBuilder(fd.IsDistinct(), fd.CaseSensitive)
		fs(f)
		o.free[key] = f
		return
	}

	// Unknown key.  Report error?
}

func (o *engineObjBuilder) Seal() *EngineObj {
	if o == nil {
		return nil
	}
	return &EngineObj{
		Source:  o.source,
		Numeric: sealDescriptorMap(o.numeric),
		Enum:    sealDescriptorMap(o.enum),
		Free:    sealDescriptorMap(o.free),
	}
}

func mutateDescriptorMap[T string | float64](
	m map[string]descriptor.ImmutableDescriptorValue[T],
) map[string]descriptor.DescriptorValueBuilder[T] {
	ret := make(map[string]descriptor.DescriptorValueBuilder[T])
	for k, v := range m {
		ret[k] = v.Copy()
	}
	return ret
}

func sealDescriptorMap[T string | float64](
	m map[string]descriptor.DescriptorValueBuilder[T],
) map[string]descriptor.ImmutableDescriptorValue[T] {
	ret := make(map[string]descriptor.ImmutableDescriptorValue[T])
	for k, v := range m {
		ret[k] = v.Seal()
	}
	return ret
}

type alterAction struct {
	val DescriptorValues
}

func (a *alterAction) addFloat(v descriptor.DescriptorValueBuilder[float64]) {
	v.AddList(a.val.Number)
}

func (a *alterAction) addString(v descriptor.DescriptorValueBuilder[string]) {
	v.AddList(a.val.Text)
}

func (a *alterAction) addDistinctFloat(v descriptor.DescriptorValueBuilder[float64]) {
	if v.IsDistinct() {
		// does it on its own
		v.AddList(a.val.Number)
		return
	}
	missing := make([]float64, 0, len(a.val.Number))
	for _, x := range a.val.Number {
		if !v.Has(x) {
			missing = append(missing, x)
		}
	}
	v.AddList(missing)
}

func (a *alterAction) addDistinctString(v descriptor.DescriptorValueBuilder[string]) {
	if v.IsDistinct() {
		// does it on its own
		v.AddList(a.val.Text)
		return
	}
	missing := make([]string, 0, len(a.val.Number))
	for _, x := range a.val.Text {
		if !v.Has(x) {
			missing = append(missing, x)
		}
	}
	v.AddList(missing)
}

func (a *alterAction) removeFloat(v descriptor.DescriptorValueBuilder[float64]) {
	for _, n := range a.val.Number {
		v.RemoveOnce(n)
	}
}

func (a *alterAction) removeString(v descriptor.DescriptorValueBuilder[string]) {
	for _, n := range a.val.Text {
		v.RemoveOnce(n)
	}
}

func (a *alterAction) removeDistinctFloat(v descriptor.DescriptorValueBuilder[float64]) {
	for _, n := range a.val.Number {
		v.RemoveAll(n)
	}
}

func (a *alterAction) removeDistinctString(v descriptor.DescriptorValueBuilder[string]) {
	for _, n := range a.val.Text {
		v.RemoveAll(n)
	}
}

func (a *alterAction) setFloat(v descriptor.DescriptorValueBuilder[float64]) {
	v.Clear()
	v.AddList(a.val.Number)
}

func (a *alterAction) setString(v descriptor.DescriptorValueBuilder[string]) {
	v.Clear()
	v.AddList(a.val.Text)
}
