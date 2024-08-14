// Under the Apache-2.0 License
package sont

import (
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/sources"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
	"github.com/groboclown/qazaar-testing/rule-engine/schema/ontology"
)

type AllowedDescriptors struct {
	Problems *problem.ProblemSet
	enums    map[string]*EnumDesc
	frees    map[string]*FreeDesc
	numerics map[string]*NumericDesc
	keyTypes map[string]DescriptorType
	sources  *sources.SourceGen
}

type DescriptorType int

const (
	UnknownDescriptorType DescriptorType = iota
	EnumDescriptorType
	FreeDescriptorType
	NumericDescriptorType
)

var keyName = map[DescriptorType]string{
	EnumDescriptorType:    "enum",
	FreeDescriptorType:    "free",
	NumericDescriptorType: "number",
}

type Descriptor interface {
	Type() DescriptorType
	KeyName() string
	IsDistinct() bool
	IsCaseSensitive() bool
}

type EnumDesc struct {
	Distinct     bool
	Enum         map[string]string
	Key          string
	MaximumCount int
	Comments     []string
	Sources      []sources.Source
}

type FreeDesc struct {
	CaseSensitive bool
	Constraints   []ValueConstraint
	Distinct      bool
	Key           string
	MaximumLength int
	MaximumCount  int
	Comments      []string
	Sources       []sources.Source
}

type NumericDesc struct {
	Distinct     bool
	Key          string
	Maximum      float64
	Minimum      float64
	MaximumCount int
	Comments     []string
	Sources      []sources.Source
}

type ValueConstraint struct {
	Format   *string
	Pattern  *string
	Type     ontology.ValueConstraintType
	Comments []string
	Sources  []sources.Source
}

// New creates a new, shared Descriptors structure.
func New() *AllowedDescriptors {
	return &AllowedDescriptors{
		keyTypes: make(map[string]DescriptorType),
		enums:    make(map[string]*EnumDesc),
		frees:    make(map[string]*FreeDesc),
		numerics: make(map[string]*NumericDesc),
		Problems: problem.New(),
		sources:  sources.SourceGenerator(),
	}
}
