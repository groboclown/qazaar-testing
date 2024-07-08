// Under the Apache-2.0 License
package validate

import (
	"context"
	"sync"

	"github.com/groboclown/qazaar-testing/rule-engine/ingest/sont"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/srule"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
)

func ValidateRuleAsync(
	rule *srule.Rule,
	ont *sont.AllowedDescriptors,
	probs problem.Adder,
	ctx context.Context,
) <-chan bool {
	ret := make(chan bool)

	go func() {
		defer func() {
			ret <- onDefer("rule", nil, probs)
			close(ret)
		}()

		var wg sync.WaitGroup

		if rule != nil {
			ValidateMatchersAsync(rule.Matchers, ont, &wg, rule.Sources, probs)
			for _, c := range rule.Conformities {
				ValidateConformityAsync(&c, ont, &wg, probs)
			}
		}

		wg.Wait()
	}()

	return ret
}

func ValidateConformityAsync(
	mat *srule.LeveledMatcher,
	ont *sont.AllowedDescriptors,
	wg *sync.WaitGroup,
	probs problem.Adder,
) {
	// TODO also need to check the level to see if the config has a reference to it.
	if mat != nil {
		ValidateMatchersAsync(mat.Matchers, ont, wg, mat.Sources, probs)
	}
}
