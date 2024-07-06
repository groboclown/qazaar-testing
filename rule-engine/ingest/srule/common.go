// Under the Apache-2.0 License
package srule

import (
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/comments"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/sources"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
	"github.com/groboclown/qazaar-testing/rule-engine/schema/rules"
)

func joinVariableMap(
	vars []rules.Variable,
	src *sources.RulesSource,
	probs *problem.ProblemSet,
) map[string]*VariableDef {
	ret := make(map[string]*VariableDef)
	for _, v := range vars {
		s := src.DocumentSources(v.Sources)
		if _, ok := ret[v.Name]; ok {
			probs.AddWarning(
				s,
				"Duplicate variable '%s'",
				v.Name,
			)
			continue
		}
		ret[v.Name] = &VariableDef{
			Name:        v.Name,
			Description: v.Description,
			Type:        v.Type,
			Comments:    comments.JoinRuleComments(v.Comment, v.Comments),
			Sources:     s,
		}
	}
	return ret
}
