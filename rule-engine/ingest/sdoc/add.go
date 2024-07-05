// Under the Apache-2.0 License
package sdoc

import (
	"github.com/groboclown/qazaar-testing/rule-engine/schema/document"
	"github.com/groboclown/qazaar-testing/rule-engine/sources"
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
) DocumentObject {
	return DocumentObject{
		Comments:    joinComments(obj),
		Descriptors: obj.Descriptors,
		Id:          obj.Id,
		Sources:     prep.DocumentObject(obj),
	}
}

func joinComments(obj *document.DocumentObject) []string {
	ret := make([]string, 0)
	if obj != nil {
		if obj.Comment != nil && *obj.Comment != "" {
			ret = append(ret, string(*obj.Comment))
		}
		for _, c := range obj.Comments {
			if c != "" {
				ret = append(ret, string(c))
			}
		}
	}
	return ret
}
