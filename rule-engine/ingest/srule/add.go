// Under the Apache-2.0 License
package srule

import (
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/comments"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/descriptor"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/sources"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
	"github.com/groboclown/qazaar-testing/rule-engine/schema/rules"
)

func (r *RuleSet) Add(obj *rules.RulesV1SchemaJson) {
	if obj == nil || r == nil {
		return
	}
	prep := r.sources.PrepareRules(&obj.CommonSourceRefs)
	for _, rule := range obj.Rules {
		r.addRule(&rule, prep)
	}
	for _, group := range obj.Groups {
		r.addGroup(&group, prep)
	}
}

func (r *RuleSet) addRule(obj *rules.Rule, src *sources.RulesSource) {
	if r == nil || obj == nil {
		return
	}
	s := src.DocumentSources(obj.Sources)
	r.Rules = append(r.Rules, &Rule{
		Comments:     comments.JoinRuleComments(obj.Comment, obj.Comments),
		Sources:      s,
		Id:           string(obj.Id),
		Variables:    joinVariableMap(obj.Variables, src, r.Problems),
		Matchers:     joinMatchers(obj.MatchingDescriptors, src, r.Problems),
		Conformities: joinConformities(obj.Conformities, src, r.Problems),
	})
}

func (r *RuleSet) addGroup(obj *rules.Group, src *sources.RulesSource) {
	if r == nil || obj == nil {
		return
	}
	s := src.DocumentSources(obj.Sources)
	r.Groups = append(r.Groups, &Group{
		Comments:        comments.JoinRuleComments(obj.Comment, obj.Comments),
		Sources:         s,
		Id:              string(obj.Id),
		Variables:       joinVariableMap(obj.Variables, src, r.Problems),
		Matchers:        joinMatchers(obj.MatchingDescriptors, src, r.Problems),
		KeySharedValues: joinKeys(obj.SharedValues),
		Alterations:     joinAlterations(obj.Alterations, src, r.Problems),
		Convergences:    joinConvergences(obj.Convergences, src, r.Problems),
	})
}

func joinKeys(keys []rules.DescriptorKey) []string {
	ret := make([]string, len(keys))
	for i, k := range keys {
		ret[i] = string(k)
	}
	return ret
}

func joinAlterations(
	alts []rules.Alteration,
	src *sources.RulesSource,
	probs *problem.ProblemSet,
) []Alteration {
	ret := make([]Alteration, 0)
	for _, a := range alts {
		s := src.DocumentSources(a.Sources)
		texts, numbers := descriptor.Join(descriptor.DecodeRuleValues(a.Values))
		ret = append(ret, Alteration{
			Key:          string(a.Key),
			Action:       toAlterationAction(a.Action, s, probs),
			TextValues:   texts,
			NumberValues: numbers,
			Comments:     comments.JoinRuleComments(a.Comment, a.Comments),
			Sources:      s,
		})
	}
	return ret
}
