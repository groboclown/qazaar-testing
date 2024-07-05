// Under the Apache-2.0 License
package sources

import "github.com/groboclown/qazaar-testing/rule-engine/schema/document"

type DocumentSource struct {
	MissingRefs []document.Id

	sg   *SourceGen
	refs map[document.Id]*document.CommonDocumentSource
}

// PrepareDocument prepares a structure for extracting universal sources from document values.
func (sg *SourceGen) PrepareDocument(cdl *document.CommonDocumentSourceList) *DocumentSource {
	ret := &DocumentSource{
		MissingRefs: make([]document.Id, 0),
		sg:          sg,
		refs:        make(map[document.Id]*document.CommonDocumentSource),
	}
	if cdl != nil {
		for _, cds := range *cdl {
			ret.refs[cds.Id] = &cds
		}
	}
	return ret
}

// DocumentObject converts the sources in the DocumentObject into the universal source value.
func (ds *DocumentSource) DocumentObject(obj *document.DocumentObject) []Source {
	if obj == nil || ds == nil {
		return nil
	}
	ret := make([]Source, 0)
	for _, s := range obj.Sources {
		if ref, ok := ds.refs[s.Ref]; ok {
			ir := ds.sg.addRef(ref.Loc, ref.Rep, ref.Ver)
			ret = append(ret, Source{
				ref: ir,
				a:   s.A,
			})
		} else {
			ds.MissingRefs = append(ds.MissingRefs, s.Ref)
		}
	}
	return ret
}
