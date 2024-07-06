// Under the Apache-2.0 License
package sources

import (
	"github.com/groboclown/qazaar-testing/rule-engine/schema/rules"
)

type RulesSource struct {
	MissingRefs []rules.Id

	sg   *SourceGen
	refs map[rules.Id]*rules.CommonDocumentSource
}

// PrepareOntology prepares a structure for extracting universal sources from document values.
func (sg *SourceGen) PrepareRules(cdl *rules.CommonDocumentSourceList) *RulesSource {
	ret := &RulesSource{
		MissingRefs: make([]rules.Id, 0),
		sg:          sg,
		refs:        make(map[rules.Id]*rules.CommonDocumentSource),
	}
	if cdl != nil {
		for _, cds := range *cdl {
			ret.refs[cds.Id] = &cds
		}
	}
	return ret
}

// DocumentSources converts the sources into the universal source value.
func (ds *RulesSource) DocumentSources(obj rules.DocumentSources) []Source {
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
