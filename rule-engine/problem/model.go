// Under the Apache-2.0 License
package problem

import "github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/sources"

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

type Adder interface {
	Done()
	Error(source string, err ...error)
	Add(p Problem)
	AddError(
		sources []sources.Source,
		format string,
		args ...any,
	)
	AddWarning(
		sources []sources.Source,
		format string,
		args ...any,
	)
	AddInfo(
		sources []sources.Source,
		format string,
		args ...any,
	)
	AddProblem(
		sources []sources.Source,
		level ProblemLevel,
		format string,
		args ...any,
	)
}
