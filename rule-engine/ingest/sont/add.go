// Under the Apache-2.0 License
package sont

import (
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/internal/sel"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/comments"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/sources"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
	"github.com/groboclown/qazaar-testing/rule-engine/schema/ontology"
	"github.com/mitchellh/mapstructure"
)

// Add adds all the descriptors in the ontology document into the descriptors structure.
func (s *AllowedDescriptors) Add(obj *ontology.OntologyV1SchemaJson) {
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
func (s *AllowedDescriptors) addDescriptor(
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

func (s *AllowedDescriptors) addEnum(obj *ontology.EnumDescriptor, src *sources.OntologySource) {
	if obj == nil || s == nil {
		return
	}
	sl := src.DocumentSources(obj.Sources)
	if s.hasDup(EnumDescriptorType, obj.Key, sl, s.Problems) {
		return
	}
	s.Enum[string(obj.Key)] = &EnumDesc{
		Comments:     comments.JoinOntComments(obj.Comment, obj.Comments),
		Distinct:     obj.Distinct,
		Enum:         enumMap(obj.Enum),
		Key:          string(obj.Key),
		MaximumCount: obj.MaximumCount,
		Sources:      sl,
	}
}

func enumMap(vals []string) map[string]string {
	ret := make(map[string]string)
	for _, v := range vals {
		ret[v] = v
	}
	return ret
}

func (s *AllowedDescriptors) addFree(obj *ontology.FreeDescriptor, src *sources.OntologySource) {
	if obj == nil || s == nil {
		return
	}
	sl := src.DocumentSources(obj.Sources)
	if s.hasDup(FreeDescriptorType, obj.Key, sl, s.Problems) {
		return
	}
	s.Free[string(obj.Key)] = &FreeDesc{
		Comments:      comments.JoinOntComments(obj.Comment, obj.Comments),
		Distinct:      obj.Distinct,
		Key:           string(obj.Key),
		CaseSensitive: obj.CaseSensitive,
		Constraints:   convertConstraints(obj.Constraints, src),
		MaximumCount:  obj.MaximumCount,
		MaximumLength: obj.MaximumLength,
		Sources:       sl,
	}
}

func (s *AllowedDescriptors) addNumeric(obj *ontology.NumericDescriptor, src *sources.OntologySource) {
	if obj == nil || s == nil {
		return
	}
	sl := src.DocumentSources(obj.Sources)
	if s.hasDup(NumericDescriptorType, obj.Key, sl, s.Problems) {
		return
	}
	s.Numeric[string(obj.Key)] = &NumericDesc{
		Comments:     comments.JoinOntComments(obj.Comment, obj.Comments),
		Distinct:     obj.Distinct,
		Key:          string(obj.Key),
		Maximum:      float64(obj.Maximum),
		Minimum:      float64(obj.Minimum),
		MaximumCount: obj.MaximumCount,
		Sources:      sl,
	}
}

func (s *AllowedDescriptors) hasDup(
	t DescriptorType,
	key ontology.DescriptorKey,
	src []sources.Source,
	p *problem.ProblemSet,
) bool {
	k := string(key)
	if _, ok := s.keyTypes[k]; ok {
		p.AddWarning(
			src,
			"%s: duplicate key (%s)",
			keyName[t],
			key,
		)
		return true
	}
	s.keyTypes[k] = t
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
		Comments: comments.JoinOntComments(con.Comment, con.Comments),
		Type:     con.Type,
		Format:   con.Format,
		Pattern:  con.Pattern,
		Sources:  src.DocumentSources(con.Sources),
	}
}
