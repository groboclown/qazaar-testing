// Under the Apache-2.0 License
package sont

import (
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/internal/sel"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
	"github.com/groboclown/qazaar-testing/rule-engine/schema/ontology"
	"github.com/groboclown/qazaar-testing/rule-engine/sources"
	"github.com/mitchellh/mapstructure"
)

// Add adds all the descriptors in the ontology document into the descriptors structure.
func (s *Descriptors) Add(obj *ontology.OntologyV1SchemaJson) {
	if obj == nil || s == nil {
		return
	}
	prep := s.sources.PrepareOntology(&obj.CommonSourceRefs)
	for _, d := range obj.Descriptors {
		s.addDescriptor(&d, prep)
	}
}

// addDescriptor add the right typed descriptor.
//
// This requires adding schema knowledge into the non-generated source, due to limitations
// in the source generator.
func (s *Descriptors) addDescriptor(
	obj *ontology.OntologyV1SchemaJsonDescriptorsElem,
	src *sources.OntologySource,
) {
	if obj == nil {
		return
	}
	err := sel.TypeSelector(
		*obj,
		"type",
		sel.SelectHandlerMap{
			"enum": func(val map[string]any) error {
				var desc ontology.EnumDescriptor
				err := mapstructure.Decode(val, &desc)
				if err != nil {
					return err
				}
				s.addEnum(&desc, src)
				return nil
			},
			"free": func(val map[string]any) error {
				var desc ontology.FreeDescriptor
				err := mapstructure.Decode(val, &desc)
				if err != nil {
					return err
				}
				s.addFree(&desc, src)
				return nil
			},
			"number": func(val map[string]any) error {
				var desc ontology.NumericDescriptor
				err := mapstructure.Decode(val, &desc)
				if err != nil {
					return err
				}
				s.addNumeric(&desc, src)
				return nil
			},
		},
	)
	if err != nil {
		s.Problems.AddError(
			nil,
			"error decoding descriptor: %s",
			err.Error(),
		)
	}
}

func (s *Descriptors) addEnum(obj *ontology.EnumDescriptor, src *sources.OntologySource) {
	if obj == nil || s == nil {
		return
	}
	sl := src.DocumentSources(obj.Sources)
	if hasDup("enum", obj.Key, s.Enum, sl, s.Problems) {
		return
	}
	s.Enum[string(obj.Key)] = &EnumDesc{
		Comments:     joinComments(obj.Comment, obj.Comments),
		Distinct:     obj.Distinct,
		Enum:         obj.Enum,
		Key:          string(obj.Key),
		MaximumCount: obj.MaximumCount,
		Sources:      sl,
	}
}

func (s *Descriptors) addFree(obj *ontology.FreeDescriptor, src *sources.OntologySource) {
	if obj == nil || s == nil {
		return
	}
	sl := src.DocumentSources(obj.Sources)
	if hasDup("free", obj.Key, s.Enum, sl, s.Problems) {
		return
	}
	s.Free[string(obj.Key)] = &FreeDesc{
		Comments:      joinComments(obj.Comment, obj.Comments),
		Distinct:      obj.Distinct,
		Key:           string(obj.Key),
		CaseSensitive: obj.CaseSensitive,
		Constraints:   convertConstraints(obj.Constraints, src),
		MaximumCount:  obj.MaximumCount,
		MaximumLength: obj.MaximumLength,
		Sources:       sl,
	}
}

func (s *Descriptors) addNumeric(obj *ontology.NumericDescriptor, src *sources.OntologySource) {
	if obj == nil || s == nil {
		return
	}
	sl := src.DocumentSources(obj.Sources)
	if hasDup("enum", obj.Key, s.Enum, sl, s.Problems) {
		return
	}
	s.Numeric[string(obj.Key)] = &NumericDesc{
		Comments:     joinComments(obj.Comment, obj.Comments),
		Distinct:     obj.Distinct,
		Key:          string(obj.Key),
		Maximum:      obj.Maximum,
		Minimum:      obj.Minimum,
		MaximumCount: obj.MaximumCount,
		Sources:      sl,
	}
}

func hasDup[T EnumDesc | FreeDesc | NumericDesc](
	name string,
	key ontology.DescriptorKey,
	m map[string]*T,
	src []sources.Source,
	p *problem.ProblemSet,
) bool {
	if _, ok := m[string(key)]; ok {
		p.AddWarning(
			src,
			"%s: duplicate key '%s'",
			name,
			key,
		)
	}
	return false
}

func convertConstraints(
	cons []ontology.ValueConstraint,
	src *sources.OntologySource,
) []ValueConstraint {
	ret := make([]ValueConstraint, len(cons))
	for i, v := range cons {
		ret[i] = convertConstraint(&v, src)
	}
	return ret
}

func convertConstraint(con *ontology.ValueConstraint, src *sources.OntologySource) ValueConstraint {
	return ValueConstraint{
		Comments: joinComments(con.Comment, con.Comments),
		Type:     con.Type,
		Format:   con.Format,
		Maximum:  con.Maximum,
		Minimum:  con.Minimum,
		Pattern:  con.Pattern,
		Sources:  src.DocumentSources(con.Sources),
	}
}

func joinComments(com *ontology.Comment, cl ontology.CommentList) []string {
	ret := make([]string, 0)
	if com != nil && *com != "" {
		ret = append(ret, string(*com))
	}
	for _, c := range cl {
		if c != "" {
			ret = append(ret, string(c))
		}
	}
	return ret
}
