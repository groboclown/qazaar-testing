// Under the Apache-2.0 License
package sources

import (
	"github.com/groboclown/qazaar-testing/rule-engine/schema/ontology"
)

type OntologySource struct {
	MissingRefs []ontology.Id

	sg   *SourceGen
	refs map[ontology.Id]*ontology.CommonDocumentSource
}

// PrepareOntology prepares a structure for extracting universal sources from document values.
func (sg *SourceGen) PrepareOntology(cdl *ontology.CommonDocumentSourceList) *OntologySource {
	ret := &OntologySource{
		MissingRefs: make([]ontology.Id, 0),
		sg:          sg,
		refs:        make(map[ontology.Id]*ontology.CommonDocumentSource),
	}
	if cdl != nil {
		for _, cds := range *cdl {
			ret.refs[cds.Id] = &cds
		}
	}
	return ret
}

// Enum converts the sources in the EnumDescriptor into the universal source value.
func (ds *OntologySource) Enum(obj *ontology.EnumDescriptor) []Source {
	if ds == nil {
		return nil
	}
	return ds.DocumentSources(obj.Sources)
}

// Free converts the sources in the FreeDescriptor into the universal source value.
func (ds *OntologySource) Free(obj *ontology.FreeDescriptor) []Source {
	if ds == nil {
		return nil
	}
	return ds.DocumentSources(obj.Sources)
}

// Numeric converts the sources in the NumericDescriptor into the universal source value.
func (ds *OntologySource) Numeric(obj *ontology.NumericDescriptor) []Source {
	if ds == nil {
		return nil
	}
	return ds.DocumentSources(obj.Sources)
}

// DocumentSources converts the sources into the universal source value.
func (ds *OntologySource) DocumentSources(obj ontology.DocumentSources) []Source {
	if obj == nil || ds == nil {
		return nil
	}
	ret := make([]Source, 0)
	for _, s := range obj {
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
