// Under the Apache-2.0 License
package matcher_test

import (
	"testing"

	"github.com/groboclown/qazaar-testing/rule-engine/engine/matcher"
	"github.com/groboclown/qazaar-testing/rule-engine/engine/obj"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/descriptor"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/sources"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/srule"
)

func Test_IsContainsMatch(t *testing.T) {
	t.Run("simple-number-all-ok", func(t *testing.T) {
		o := mkObj().numeric("e1", false, 1.5)
		c := &srule.ContainsMatcher{
			Operation: srule.ContainsAll,
			Count:     false,
			Distinct:  false,
			Key:       "e1",
			Checks: srule.ValueCheckSet{
				Numeric: []srule.NumericBoundsCheck{{Min: 1, Max: 2}},
				// Text: []srule.StringCheck{{R: regexp.MustCompile("^a+$")}},
			},
		}
		if !matcher.IsContainsMatch(o.eo, c) {
			t.Error("did not match")
		}
	})
	t.Run("simple-number-all-fail", func(t *testing.T) {
		o := mkObj().numeric("e1", false, 1.5)
		c := &srule.ContainsMatcher{
			Operation: srule.ContainsAll,
			Count:     false,
			Distinct:  false,
			Key:       "e1",
			Checks: srule.ValueCheckSet{
				Numeric: []srule.NumericBoundsCheck{{Min: 2, Max: 3}},
				// Text: []srule.StringCheck{{R: regexp.MustCompile("^a+$")}},
			},
		}
		if matcher.IsContainsMatch(o.eo, c) {
			t.Error("incorrectly matched")
		}
	})
}

type objBuilder struct {
	eo *obj.EngineObj
}

func mkObj() *objBuilder {
	return &objBuilder{&obj.EngineObj{
		Source: obj.ObjSource{
			Parents:   nil,
			Construct: nil,
			Source:    []sources.Source{{}},
		},
		Numeric: make(map[string]descriptor.DescriptorValueBuilder[float64]),
		Enum:    make(map[string]descriptor.DescriptorValueBuilder[string]),
		Free:    make(map[string]descriptor.DescriptorValueBuilder[string]),
	}}
}

func (o *objBuilder) enum(key string, distinct bool, val ...string) *objBuilder {
	b := descriptor.NewTextBuilder(distinct, true)
	b.AddList(val)
	o.eo.Enum[key] = b
	return o
}

func (o *objBuilder) free(key string, distinct bool, caseSensitive bool, val ...string) *objBuilder {
	b := descriptor.NewTextBuilder(distinct, caseSensitive)
	b.AddList(val)
	o.eo.Free[key] = b
	return o
}

func (o *objBuilder) numeric(key string, distinct bool, val ...float64) *objBuilder {
	b := descriptor.NewNumericBuilder(distinct)
	b.AddList(val)
	o.eo.Numeric[key] = b
	return o
}
