// Under the Apache-2.0 License
package problem

import (
	"fmt"

	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/sources"
)

func (s *ProblemSet) Done() {
	// do nothing
}

// Add adds a problem to the set.
func (s *ProblemSet) Add(p Problem) {
	if s == nil {
		return
	}
	// TODO could enforce a uniqueness here.
	s.p = append(s.p, p)
}

// Merges another problem set into this one.
func (s *ProblemSet) Merge(o *ProblemSet) {
	if s == nil || o == nil {
		return
	}
	for _, p := range o.p {
		s.Add(p)
	}
}

func (ps *ProblemSet) AddError(
	sources []sources.Source,
	format string,
	args ...any,
) {
	ps.AddProblem(sources, Err, format, args...)
}

func (ps *ProblemSet) AddWarning(
	sources []sources.Source,
	format string,
	args ...any,
) {
	ps.AddProblem(sources, Warn, format, args...)
}

func (ps *ProblemSet) AddInfo(
	sources []sources.Source,
	format string,
	args ...any,
) {
	ps.AddProblem(sources, Info, format, args...)
}

func (ps *ProblemSet) AddProblem(
	sources []sources.Source,
	level ProblemLevel,
	format string,
	args ...any,
) {
	ps.Add(Problem{
		Level:   level,
		Message: fmt.Sprintf(format, args...),
		Sources: sources,
	})
}
