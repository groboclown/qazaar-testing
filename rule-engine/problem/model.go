// Under the Apache-2.0 License
package problem

import "github.com/groboclown/qazaar-testing/rule-engine/sources"

type ProblemLevel int

const (
	Quiet ProblemLevel = iota
	Info
	Warn
	Err
)

// Problem represents a single problem.
type Problem struct {
	Level   ProblemLevel
	Message string
	Sources []sources.Source
}

// ProblemSet contains many problems.
type ProblemSet struct {
	p []Problem
}

func New() *ProblemSet {
	return &ProblemSet{p: make([]Problem, 0)}
}
