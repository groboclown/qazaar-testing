// Under the Apache-2.0 License
package sdoc

import (
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/comments"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/descriptor"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/sources"
	"github.com/groboclown/qazaar-testing/rule-engine/schema/document"
)

// Add adds in the documents from the data-exchange format into the simplified form.
func (d *Documents) Add(src *document.DocumentDescriptionV1SchemaJson) {
	if src == nil || d == nil {
		return
	}
	prep := d.sources.PrepareDocument(&src.CommonSourceRefs)

	for _, obj := range src.Objects {
		d.Objects = append(d.Objects, d.updateSources(&obj, prep))
	}
}

func (d *Documents) updateSources(
	obj *document.DocumentObject,
	prep *sources.DocumentSource,
) *DocumentObject {
	return &DocumentObject{
		Comments:    comments.JoinDocComments(obj),
		Descriptors: descriptor.JoinDocumentDescriptors(obj.Descriptors),
		Id:          obj.Id,
		Sources:     prep.DocumentObject(obj),
	}
}
