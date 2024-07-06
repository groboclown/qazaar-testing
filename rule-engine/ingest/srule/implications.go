// Under the Apache-2.0 License
package srule

import (
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/sources"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
	"github.com/groboclown/qazaar-testing/rule-engine/schema/rules"
)

func joinConformities(
	conf []rules.ConformityImplication,
	src *sources.RulesSource,
	probs *problem.ProblemSet,
) []LeveledMatcher {
	byLevel := make(map[rules.ImplicationLevel]*LeveledMatcher)
	for _, c := range conf {
		m, ok := byLevel[c.Level]
		if !ok {
			m = &LeveledMatcher{
				Level: string(c.Level),
				Matchers: &MatchingDescriptorSet{
					Collection: make([]CollectionMatcher, 0),
					Contains:   make([]ContainsMatcher, 0),
				},
			}
			byLevel[c.Level] = m
		}
		addConformity(m, &c, src, probs)
	}

	ret := make([]LeveledMatcher, len(byLevel))
	i := 0
	for _, v := range byLevel {
		ret[i] = *v
		i++
	}
	return ret
}

func addConformity(
	m *LeveledMatcher,
	conf *rules.ConformityImplication,
	src *sources.RulesSource,
	probs *problem.ProblemSet,
) {
	if m == nil || conf == nil {
		return
	}
	addMatcher(m.Matchers, conf, src, probs)
}

func joinConvergences(
	conv []rules.ConvergenceImplication,
	src *sources.RulesSource,
	probs *problem.ProblemSet,
) []Convergence {
	ret := make([]Convergence, 0)
	for _, c := range conv {
		s := src.DocumentSources(c.Sources)
		ret = append(ret, Convergence{
			Key:      string(c.Key),
			Level:    string(c.Level),
			Distinct: c.Distinct,
			Requires: toConvergenceType(c.Requires, s, probs),
		})
	}
	return ret
}
