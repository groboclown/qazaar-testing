// Under the Apache-2.0 License
package runner

import (
	"sync"

	"github.com/groboclown/qazaar-testing/rule-engine/engine/matcher"
	"github.com/groboclown/qazaar-testing/rule-engine/engine/obj"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/srule"
)

func checkAllAgainstRules(
	all []*obj.EngineObj,
	rules []*srule.Rule,
) []*RuleProblem {
	ret := make([]*RuleProblem, 0)
	for p := range checkAllAgainstRulesAsync(all, rules) {
		ret = append(ret, p)
	}
	return ret
}

func checkAllAgainstRulesAsync(
	all []*obj.EngineObj,
	rules []*srule.Rule,
) <-chan *RuleProblem {
	ret := make(chan *RuleProblem)
	go func() {
		defer close(ret)

		var wg sync.WaitGroup
		for _, o := range all {
			for _, r := range rules {
				wg.Add(1)
				go func() {
					defer wg.Done()
					checkAgainstRule(o, r, ret)
				}()
			}
		}
		wg.Wait()
	}()
	return ret
}

// checkAgainstRule validates the object against the rule, and adds violations to the problem set.
func checkAgainstRule(
	o *obj.EngineObj,
	rule *srule.Rule,
	probs chan<- *RuleProblem,
) {
	if matcher.IsMatch(o, rule.Matchers) {
		// The object must conform.
		for _, c := range rule.Conformities {
			if !matcher.IsMatch(o, c.Matchers) {
				probs <- &RuleProblem{
					obj:     o,
					rule:    rule,
					matcher: &c,
				}
			}
		}
	}
}
