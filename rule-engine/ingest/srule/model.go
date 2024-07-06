// Under the Apache-2.0 License
package srule

import (
	"regexp"

	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/sources"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
)

// RuleSet contains a collection of all processed rules and groups.
type RuleSet struct {
	Rules    []*Rule
	Groups   []*Group
	Problems *problem.ProblemSet
	sources  *sources.SourceGen
}

type Rule struct {
	Id           string
	Variables    map[string]*VariableDef
	Matchers     *MatchingDescriptorSet
	Conformities []LeveledMatcher
	Comments     []string
	Sources      []sources.Source
}

type Group struct {
	Id        string
	Variables map[string]*VariableDef
	Matchers  *MatchingDescriptorSet
	// SharedValues
	Alterations  []Alteration
	Convergences []Convergence
	Comments     []string
	Sources      []sources.Source
}

type ConvergenceType int

const (
	AllMatch ConvergenceType = iota
	Disjoint
)

type Convergence struct {
	Key      string
	Level    string
	Distinct bool
	Requires ConvergenceType
}

type AlterationAction int

const (
	AddAction AlterationAction = iota
	AddDistinctAction
	RemoveAction
	RemoveDistinctAction
	SetAction
)

type Alteration struct {
	Key          string
	Action       AlterationAction
	TextValues   []string
	NumberValues []float64
}

type LeveledMatcher struct {
	Level    string
	Matchers *MatchingDescriptorSet
}

type MatchingDescriptorSet struct {
	Collection []CollectionMatcher
	Contains   []ContainsMatcher
}

type CollectionOperation int

const (
	AndCollection CollectionOperation = iota
	OrCollection
	NotCollection
)

type CollectionMatcher struct {
	Operation CollectionOperation
	Matchers  *MatchingDescriptorSet
}

type ContainsOperation int

const (
	ContainsAll ContainsOperation = iota
	ContainsSome
	ContainsOnly
	ContainsExactly
)

type ContainsMatcher struct {
	Operation ContainsOperation
	Count     bool
	Distinct  bool
	Key       string
	Checks    ValueCheckSet
}

type ValueCheckSet struct {
	Collection []CollectionCheck
	String     []StringCheck
	Numeric    []NumericBoundsCheck
}

type CollectionCheckOperation int

const (
	OrCheck CollectionCheckOperation = iota
	AndCheck
	NotCheck
)

type CollectionCheck struct {
	Operation  CollectionCheckOperation
	Collection ValueCheckSet
}

type StringCheck *regexp.Regexp

type NumericBoundsCheck struct {
	Min float64
	Max float64
}

type VariableDef struct {
	Comments    []string
	Description *string
	Name        string
	Type        string
	Sources     []sources.Source
}

func New() *RuleSet {
	return &RuleSet{
		Rules:    make([]*Rule, 0),
		Groups:   make([]*Group, 0),
		Problems: problem.New(),
		sources:  sources.SourceGenerator(),
	}
}
