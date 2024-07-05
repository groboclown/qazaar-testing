// Under the Apache-2.0 License
package sont

import (
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
	"github.com/groboclown/qazaar-testing/rule-engine/schema/ontology"
	"github.com/groboclown/qazaar-testing/rule-engine/sources"
)

type Descriptors struct {
	Enum     map[string]*EnumDesc
	Free     map[string]*FreeDesc
	Numeric  map[string]*NumericDesc
	Problems *problem.ProblemSet
	sources  *sources.SourceGen
}

type EnumDesc struct {
	Comments     []string
	Distinct     bool
	Enum         []string
	Key          string
	MaximumCount int
	Sources      []sources.Source
}

type FreeDesc struct {
	Comments      []string
	CaseSensitive bool
	Constraints   []ValueConstraint
	Distinct      bool
	Key           string
	MaximumCount  int
	MaximumLength int
	Sources       []sources.Source
}

type NumericDesc struct {
	Comments     []string
	Distinct     bool
	Key          string
	Maximum      ontology.DescriptorNumericValue
	MaximumCount int
	Minimum      ontology.DescriptorNumericValue
	Sources      []sources.Source
}

type ValueConstraint struct {
	Comments []string
	Format   *string
	Maximum  *ontology.DescriptorNumericValue
	Minimum  *ontology.DescriptorNumericValue
	Pattern  *string
	Sources  []sources.Source
	Type     ontology.ValueConstraintType
}

// New creates a new, shared Descriptors structure.
func New() *Descriptors {
	return &Descriptors{
		Enum:     make(map[string]*EnumDesc),
		Free:     make(map[string]*FreeDesc),
		Numeric:  make(map[string]*NumericDesc),
		Problems: problem.New(),
		sources:  sources.SourceGenerator(),
	}
}
