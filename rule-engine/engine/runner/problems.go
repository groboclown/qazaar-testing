// Under the Apache-2.0 License
package runner

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/groboclown/qazaar-testing/rule-engine/engine/obj"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/sources"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/srule"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
)

type RuleProblem struct {
	obj     *obj.EngineObj
	rule    *srule.Rule
	matcher *srule.LeveledMatcher
}

type MemberValues struct {
	Members []*obj.EngineObj
	Text    []string
	Number  []float64
}

type ConvProblem struct {
	GroupId    string
	Mismatched []MemberValues
	Conv       *srule.Convergence
}

func addRuleProblems(
	adder problem.Adder,
	levelMap map[string]problem.ProblemLevel,
	probs []*RuleProblem,
) {
	for _, p := range probs {
		adder.Add(ruleAsProblem(levelMap, p))
	}
}

func ruleAsProblem(
	levelMap map[string]problem.ProblemLevel,
	prob *RuleProblem,
) problem.Problem {
	sources := make([]sources.Source, 0)
	sources = append(sources, prob.obj.Source.Source...)
	sources = append(sources, prob.rule.Sources...)
	sources = append(sources, prob.matcher.Sources...)

	return problem.Problem{
		Level:   errLevel(prob.matcher.Level, levelMap),
		Message: ruleProblemMessage(prob),
		Sources: sources,
		Context: prob,
	}
}

func ruleProblemMessage(prob *RuleProblem) string {
	return fmt.Sprintf(
		"Rule %s violation (%s) for %s: %s",
		prob.rule.Id,
		prob.matcher.Level,
		prob.obj.String(),
		obj.ObjSource{Source: prob.matcher.Sources}.String(), // Comments instead?
	)
}

func errLevel(
	level string,
	levelMap map[string]problem.ProblemLevel,
) problem.ProblemLevel {
	if v, ok := levelMap[level]; ok {
		return v
	}
	return problem.Err
}

func convAsProblem(
	levelMap map[string]problem.ProblemLevel,
	prob *ConvProblem,
) problem.Problem {
	sources := make([]sources.Source, 0)
	for _, m := range prob.Mismatched {
		for _, o := range m.Members {
			sources = append(sources, o.Source.Source...)
		}
	}
	sources = append(sources, prob.Conv.Sources...)

	return problem.Problem{
		Level:   errLevel(prob.Conv.Level, levelMap),
		Message: convProblemMessage(prob),
		Sources: sources,
		Context: prob,
	}
}

func convProblemMessage(prob *ConvProblem) string {
	groups := make([]string, len(prob.Mismatched))
	for i, o := range prob.Mismatched {
		groups[i] = o.String()
	}

	return fmt.Sprintf(
		"Group %s: Convergence %s violation (%s); %s",
		prob.GroupId,
		prob.Conv.Key,
		prob.Conv.Level,
		strings.Join(groups, ", "),
	)
}

func (m MemberValues) String() string {
	values := make([]string, 0, len(m.Number)+len(m.Text))
	for _, v := range m.Number {
		values = append(values, strconv.FormatFloat(v, 'f', 4, 64))
	}
	values = append(values, m.Text...)
	objs := make([]string, len(m.Members))
	for i, o := range m.Members {
		objs[i] = o.String()
	}
	return fmt.Sprintf("%s (contain %s)", strings.Join(objs, ", "), strings.Join(values, ", "))
}
