// Under the Apache-2.0 License
package obj

import (
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/sdoc"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/descriptor"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/sources"
)

// ObjSource traces the source of the object back to its representative origins.
//
// Synthetic (e.g. constructed) will always have a `Construct` value, referencing
// the group rule ID that created this object.  If the
// synthetic construction comes from other objects, the `Parents` value will
// contain those parents sources.  Note that any parent may, itself, have a
// synthetic source.
type ObjSource struct {
	Parents   []*ObjSource
	Construct *string
	Source    []sources.Source
}

// DescriptorValue is an engine-specific representation of the value for any key.
//
// For numeric descriptors, the Text will be nil.  For string descriptors, the Number
// will be nil.
type DescriptorValue struct {
	Text   []string
	Number []float64
}

// EngineObj refers to an object participating in the engine rule evaluation.
//
// These may come either directly from a source, or the engine may construct it
// based on rules.
type EngineObj struct {
	Source ObjSource

	// Note: though these are modifiable values, treat them as read-only once
	//       returned.
	Numeric map[string]descriptor.DescriptorValueBuilder[float64]
	Enum    map[string]descriptor.DescriptorValueBuilder[string]
	Free    map[string]descriptor.DescriptorValueBuilder[string]
}

// ObjFactory constructs new engine objects from sources.
type ObjFactory interface {
	// FromDocument creates the object from a document description object.
	FromDocument(doc *sdoc.DocumentObject) *EngineObj

	// FromGroup creates a synthetic object based on a group rule.
	//
	// The group rule may need to perform alterations upon the generated object.
	FromGroup(members []*EngineObj, groupSrc string) *EngineObj
}

func newObj(
	parents []*ObjSource,
	construct *string,
	source []sources.Source,
) *EngineObj {
	return &EngineObj{
		Source: ObjSource{
			Parents:   parents,
			Construct: construct,
			Source:    source,
		},
		Numeric: make(map[string]descriptor.DescriptorValueBuilder[float64]),
		Enum:    make(map[string]descriptor.DescriptorValueBuilder[string]),
		Free:    make(map[string]descriptor.DescriptorValueBuilder[string]),
	}
}
