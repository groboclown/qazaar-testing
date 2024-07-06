// Under the Apache-2.0 License
package srule

import (
	"regexp"

	"github.com/groboclown/qazaar-testing/rule-engine/ingest/internal/sel"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/sources"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
	"github.com/groboclown/qazaar-testing/rule-engine/schema/rules"
	"github.com/mitchellh/mapstructure"
)

func joinMatchers[
	T rules.RuleMatchingDescriptorsElem |
		rules.GroupMatchingDescriptorsElem |
		rules.NotMatcherMatcher |
		any](
	matchers []T,
	src *sources.RulesSource,
	probs *problem.ProblemSet,
) *MatchingDescriptorSet {
	ret := MatchingDescriptorSet{
		Collection: make([]CollectionMatcher, 0),
		Contains:   make([]ContainsMatcher, 0),
	}
	for _, m := range matchers {
		addMatcher(&ret, &m, src, probs)
	}
	return &ret
}

func addMatcher[
	T rules.RuleMatchingDescriptorsElem |
		rules.GroupMatchingDescriptorsElem |
		rules.ConformityImplicationMatcher | any](
	m *MatchingDescriptorSet,
	obj *T,
	src *sources.RulesSource,
	probs *problem.ProblemSet,
) {
	if obj == nil || m == nil {
		return
	}
	err := sel.TypeSelector(
		*obj,
		"type",
		sel.SelectHandlerMap{
			string(rules.CollectionMatcherTypeOr): func(val map[string]any) error {
				return newCollectionMatcher(
					m, OrCollection, val, src, probs,
				)
			},
			string(rules.CollectionMatcherTypeAnd): func(val map[string]any) error {
				return newCollectionMatcher(
					m, AndCollection, val, src, probs,
				)
			},
			string(rules.NotMatcherTypeNot): func(val map[string]any) error {
				var match rules.NotMatcher
				err := mapstructure.Decode(val, &match)
				if err != nil {
					return err
				}
				m.Collection = append(m.Collection, CollectionMatcher{
					Operation: NotCollection,
					Matchers:  joinMatchers([]rules.NotMatcherMatcher{match.Matcher}, src, probs),
				})
				return nil
			},
			string(rules.ContainsMatcherTypeContainsSome): func(val map[string]any) error {
				return newContainsMatcher(
					m, ContainsSome, val, src, probs,
				)
			},
			string(rules.ContainsMatcherTypeContainsAll): func(val map[string]any) error {
				return newContainsMatcher(
					m, ContainsAll, val, src, probs,
				)
			},
			string(rules.ContainsMatcherTypeContainsExactly): func(val map[string]any) error {
				return newContainsMatcher(
					m, ContainsExactly, val, src, probs,
				)
			},
			string(rules.ContainsMatcherTypeContainsOnly): func(val map[string]any) error {
				return newContainsMatcher(
					m, ContainsOnly, val, src, probs,
				)
			},
		},
	)
	if err != nil {
		probs.AddError(
			nil,
			"error decoding matcher: %s",
			err.Error(),
		)
	}
}

func newCollectionMatcher(
	m *MatchingDescriptorSet,
	operation CollectionOperation,
	val map[string]any,
	src *sources.RulesSource,
	probs *problem.ProblemSet,
) error {
	var match rules.CollectionMatcher
	err := mapstructure.Decode(val, &match)
	if err != nil {
		return err
	}
	m.Collection = append(m.Collection, CollectionMatcher{
		Operation: operation,
		Matchers:  joinMatchers(match.Collection, src, probs),
	})
	return nil
}

func newContainsMatcher(
	m *MatchingDescriptorSet,
	operation ContainsOperation,
	val map[string]any,
	src *sources.RulesSource,
	probs *problem.ProblemSet,
) error {
	var match rules.ContainsMatcher
	err := mapstructure.Decode(val, match)
	if err != nil {
		return err
	}
	m.Contains = append(m.Contains, ContainsMatcher{
		Operation: operation,
		Count:     match.Count,
		Distinct:  match.Distinct,
		Key:       string(match.Key),
		Checks:    joinChecks(match.Matcher, src, probs),
	})
	return nil
}

func joinChecks(
	checks rules.ValueCheckList,
	src *sources.RulesSource,
	probs *problem.ProblemSet,
) ValueCheckSet {
	ret := &ValueCheckSet{
		Collection: make([]CollectionCheck, 0),
		String:     make([]StringCheck, 0),
		Numeric:    make([]NumericBoundsCheck, 0),
	}
	for _, c := range checks {
		err := sel.TypeSelector(
			c, "type", sel.SelectHandlerMap{
				string(rules.StringCheckTypeEqual): func(val map[string]any) error {
					return addStringCheck(val, false, ret)
				},
				string(rules.StringCheckTypePattern): func(val map[string]any) error {
					return addStringCheck(val, true, ret)
				},
				string(rules.CollectionCheckTypeAnd): func(val map[string]any) error {
					return addCollectionCheck(val, AndCheck, ret, src, probs)
				},
				string(rules.CollectionCheckTypeOr): func(val map[string]any) error {
					return addCollectionCheck(val, OrCheck, ret, src, probs)
				},
				string(rules.NotCheckTypeNot): func(val map[string]any) error {
					var check rules.NotCheck
					if err := mapstructure.Decode(val, &check); err != nil {
						return err
					}
					ret.Collection = append(ret.Collection, CollectionCheck{
						Operation:  NotCheck,
						Collection: joinChecks(rules.ValueCheckList{check.Check}, src, probs),
					})
					return nil
				},
			},
		)
		if err != nil {
			probs.AddError(
				nil,
				"error decoding value check: %s",
				err.Error(),
			)
		}
	}
	return *ret
}

func addStringCheck(val map[string]any, asRe bool, checks *ValueCheckSet) error {
	var check rules.StringCheck
	if err := mapstructure.Decode(val, &check); err != nil {
		return err
	}

	var re *regexp.Regexp = nil
	var err error
	if asRe {
		re, err = regexp.Compile(string(check.Text))
	} else {
		re, err = regexp.Compile(regexp.QuoteMeta(string(check.Text)))
	}
	if err != nil {
		return err
	}

	checks.String = append(checks.String, re)
	return nil
}

func addCollectionCheck(
	val map[string]any,
	opr CollectionCheckOperation,
	checks *ValueCheckSet,
	src *sources.RulesSource,
	probs *problem.ProblemSet,
) error {
	var check rules.CollectionCheck
	if err := mapstructure.Decode(val, &check); err != nil {
		return err
	}
	col := joinChecks(check.Collection, src, probs)
	checks.Collection = append(checks.Collection, CollectionCheck{
		Operation:  opr,
		Collection: col,
	})
	return nil
}
