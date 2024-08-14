// Under the Apache-2.0 License
package obj

import (
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/sdoc"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/descriptor"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/sources"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/sont"
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

// DescriptorValues is an engine-specific representation of the value for any key.
//
// For numeric descriptors, the Text will be nil.  For string descriptors, the Number
// will be nil.
type DescriptorValues struct {
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
	Numeric map[string]descriptor.ImmutableDescriptorValue[float64]
	Enum    map[string]descriptor.ImmutableDescriptorValue[string]
	Free    map[string]descriptor.ImmutableDescriptorValue[string]

	ont *sont.AllowedDescriptors
}

type EngineObjBuilder interface {
	Add(key string, val DescriptorValues)
	AddDistinct(key string, val DescriptorValues)
	Remove(key string, val DescriptorValues)
	RemoveDistinct(key string, val DescriptorValues)
	Set(key string, val DescriptorValues)
	Seal() *EngineObj
}

// ObjFactory constructs new engine objects from sources.
type ObjFactory interface {
	// FromDocument creates the object from a document description object.
	FromDocument(doc *sdoc.DocumentObject) *EngineObj

	// FromGroup creates a synthetic object based on a group rule.
	//
	// The group rule may need to perform alterations upon the generated object.
	FromGroup(members []*EngineObj, groupSrc string) *EngineObj

	// Empty creates an object builder.
	Empty(source ObjSource) EngineObjBuilder
}
