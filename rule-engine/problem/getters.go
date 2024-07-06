// Under the Apache-2.0 License
package problem

func (p Problem) String() string {
	return p.Message
}

func (ps *ProblemSet) HasProblems() bool {
	return len(ps.p) > 0
}

func (ps *ProblemSet) Errors() []Problem {
	return ps.ProblemsAt(Err)
}

func (ps *ProblemSet) ProblemsAt(level ProblemLevel) []Problem {
	if ps == nil {
		return nil
	}
	ret := make([]Problem, 0)
	for _, p := range ps.p {
		if p.Level == level {
			ret = append(ret, p)
		}
	}
	return ret
}

func (ps *ProblemSet) Problems() []Problem {
	return ps.p
}
