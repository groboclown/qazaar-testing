// Under the Apache-2.0 License
//
// Enum String to number conversions
package srule

import (
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/sources"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
	"github.com/groboclown/qazaar-testing/rule-engine/schema/rules"
)

func toConvergenceType(
	c rules.ConvergenceImplicationRequires,
	s []sources.Source,
	p *problem.ProblemSet,
) ConvergenceType {
	switch c {
	case rules.ConvergenceImplicationRequiresAllMatch:
		return AllMatch
	case rules.ConvergenceImplicationRequiresDisjoint:
		return Disjoint
	}
	p.AddError(
		s,
		"unsupported convergence type (%s)",
		c,
	)
	return AllMatch
}

func toAlterationAction(
	c rules.AlterationAction,
	s []sources.Source,
	p *problem.ProblemSet,
) AlterationAction {
	switch c {
	case rules.AlterationActionAdd:
		return AddAction
	case rules.AlterationActionAddDistinct:
		return AddDistinctAction
	case rules.AlterationActionRemove:
		return RemoveAction
	case rules.AlterationActionRemoveDistinct:
		return RemoveDistinctAction
	case rules.AlterationActionSet:
		return SetAction
	}
	p.AddError(
		s,
		"unsupported alteration action (%s)",
		c,
	)
	return SetAction
}
