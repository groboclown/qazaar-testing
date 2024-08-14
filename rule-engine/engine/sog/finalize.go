// Under the Apache-2.0 License
package sog

import (
	"github.com/groboclown/qazaar-testing/rule-engine/engine/obj"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/srule"
)

// asFinalizedObj applies the alterations to the object, and returns the updated object.
func asFinalizedObj(
	o *obj.EngineObj,
	alterations []srule.Alteration,
) *obj.EngineObj {
	if o == nil {
		return nil
	}

	ret := o.Alter()
	for _, a := range alterations {
		applyAlteration(ret, &a)
	}
	return ret.Seal()
}

func applyAlteration(o obj.EngineObjBuilder, alt *srule.Alteration) {
	if o == nil || alt == nil {
		return
	}

	values := obj.DescriptorValues{Text: alt.TextValues, Number: alt.NumberValues}

	switch alt.Action {
	// This may also need to take case sensitivity into account...
	case srule.AddDistinctAction:
		o.AddDistinct(alt.Key, values)
	case srule.AddAction:
		o.Add(alt.Key, values)
	case srule.RemoveDistinctAction:
		o.RemoveDistinct(alt.Key, values)
	case srule.RemoveAction:
		o.Remove(alt.Key, values)
	case srule.SetAction:
		o.Set(alt.Key, values)
	}
}
