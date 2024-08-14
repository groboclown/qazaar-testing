// Under the Apache-2.0 License
package obj

import (
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/sdoc"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/descriptor"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/sont"
)

type objFactory struct {
	ont *sont.AllowedDescriptors
}

// NewObjFactory creates a new object factory, which allows creating new engine objects.
func NewObjFactory(ont *sont.AllowedDescriptors) ObjFactory {
	return &objFactory{
		ont: ont,
	}
}

// FromDocument creates a new engine object from a document description object.
func (f *objFactory) FromDocument(doc *sdoc.DocumentObject) *EngineObj {
	if f == nil || doc == nil {
		return nil
	}
	ret := newObjBuilder(nil, nil, doc.Sources, f.ont)
	for _, d := range doc.Descriptors {
		if d == nil {
			continue
		}
		ot := f.ont.Type(d.Key)
		switch ot {
		case sont.EnumDescriptorType:
			e := f.ont.Enum(d.Key)
			if e == nil {
				// panic?
				continue
			}
			v := descriptor.NewTextBuilder(e.Distinct, true)
			ret.enum[d.Key] = v
			v.AddList(d.Text)
		case sont.FreeDescriptorType:
			e := f.ont.Free(d.Key)
			if e == nil {
				// panic?
				continue
			}
			v := descriptor.NewTextBuilder(e.Distinct, e.CaseSensitive)
			ret.free[d.Key] = v
			v.AddList(d.Text)
		case sont.NumericDescriptorType:
			e := f.ont.Numeric(d.Key)
			if e == nil {
				// panic?
				continue
			}
			v := descriptor.NewNumericBuilder(e.Distinct)
			ret.numeric[d.Key] = v
			v.AddList(d.Number)
		}
	}
	return ret.Seal()
}

// FromGroup creates a new engine object for a self-organizing group.
//
// `groupSrc` refers to the identifier for the SOG rule that defines the group.
func (f *objFactory) FromGroup(members []*EngineObj, groupSrc string) *EngineObj {
	if f == nil || members == nil || len(members) <= 0 {
		return nil
	}

	ret := newObjBuilder(copySrc(members), &groupSrc, nil, f.ont)
	for _, m := range members {
		for k, vs := range m.Enum {
			appendBuilder(k, ret.enum, vs)
		}
		for k, vs := range m.Free {
			appendBuilder(k, ret.free, vs)
		}
		for k, vs := range m.Numeric {
			appendBuilder(k, ret.numeric, vs)
		}
	}
	return ret.Seal()
}

func (f *objFactory) Empty(source ObjSource) EngineObjBuilder {
	return &engineObjBuilder{
		source:  source,
		ont:     f.ont,
		numeric: make(map[string]descriptor.DescriptorValueBuilder[float64]),
		enum:    make(map[string]descriptor.DescriptorValueBuilder[string]),
		free:    make(map[string]descriptor.DescriptorValueBuilder[string]),
	}
}

func copySrc(p []*EngineObj) []*ObjSource {
	ret := make([]*ObjSource, len(p))
	for i, o := range p {
		ret[i] = &o.Source
	}
	return ret
}

func appendBuilder[T descriptor.DescriptorValueTypes](
	key string,
	m map[string]descriptor.DescriptorValueBuilder[T],
	add descriptor.ImmutableDescriptorValue[T],
) {
	v, ok := m[key]
	if !ok {
		v = add.Copy()
		m[key] = v
	}
	v.Add(v)
}
