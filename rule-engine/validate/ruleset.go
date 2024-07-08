// Under the Apache-2.0 License
package validate

import (
	"context"
	"sync"

	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/sources"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/sont"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/srule"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
)

func ValidateRuleSetAsync(
	ruleSet *srule.RuleSet,
	ont *sont.AllowedDescriptors,
	probs problem.Adder,
	ctx context.Context,
) <-chan bool {
	ret := make(chan bool)

	go func() {
		defer func() {
			ret <- onDefer("rule set", nil, probs)
			close(ret)
		}()

		var wg sync.WaitGroup

		for _, r := range ruleSet.Rules {
			ch := ValidateRuleAsync(r, ont, probs, ctx)
			wg.Add(1)
			go func() {
				defer onDefer("rule", &wg, probs)
				<-ch
			}()
		}
		for _, g := range ruleSet.Groups {
			ch := ValidateGroupAsync(g, ont, probs, ctx)
			wg.Add(1)
			go func() {
				defer onDefer("group", &wg, probs)
				<-ch
			}()
		}

		wg.Wait()
	}()

	return ret
}

func ValidateMatchersAsync(
	mat *srule.MatchingDescriptorSet,
	ont *sont.AllowedDescriptors,
	wg *sync.WaitGroup,
	src []sources.Source,
	probs problem.Adder,
) {
	if mat == nil || ont == nil {
		return
	}
	for _, m := range mat.Collection {
		wg.Add(1)
		go func() {
			defer onDefer("collection matcher", wg, probs)
			ValidateCollectionMatcherAsync(&m, ont, wg, src, probs)
		}()
	}
	for _, m := range mat.Contains {
		wg.Add(1)
		go func() {
			defer onDefer("contains matcher", wg, probs)
			ValidateContainsMatcher(&m, ont, src, probs)
		}()
	}
}

func ValidateCollectionMatcherAsync(
	col *srule.CollectionMatcher,
	ont *sont.AllowedDescriptors,
	wg *sync.WaitGroup,
	src []sources.Source,
	probs problem.Adder,
) {
	if col == nil {
		return
	}
	ValidateMatchersAsync(col.Matchers, ont, wg, src, probs)
}

func ValidateContainsMatcher(
	con *srule.ContainsMatcher,
	ont *sont.AllowedDescriptors,
	src []sources.Source,
	probs problem.Adder,
) {
	if con == nil {
		return
	}
	typed := checkKey("contains", con.Key, ont, src, probs)
	if typed == nil {
		return
	}
	if typed.Enum != nil {
		// Numeric are allowed on count types.
		if len(con.Checks.Numeric) != 0 {
			if !con.Count {
				probs.AddError(
					src,
					"%s: enum-based contains matcher cannot have numeric values",
					con.Key,
				)
			}
		}
		if len(con.Checks.Text) != 0 {
			if con.Count {
				probs.AddError(
					src,
					"%s: enum-based contains matcher on 'count' cannot have text values",
					con.Key,
				)
			}
			// Should only be exact equality checks.
			// Should only have values in the enum.
			// However, this requires additional, captured information.
		}
	}
	if typed.Free != nil {
		if len(con.Checks.Numeric) != 0 {
			if !con.Count {
				probs.AddError(
					src,
					// dude!
					"%s: free-based contains matcher cannot have numeric values",
					con.Key,
				)
			}
		}
		if len(con.Checks.Text) != 0 {
			if con.Count {
				probs.AddError(
					src,
					"%s: free-based contains matcher on 'count' cannot have text values",
					con.Key,
				)
			}
		}
	}
	if typed.Numeric != nil {
		if len(con.Checks.Text) != 0 {
			probs.AddError(
				src,
				"%s: numeric-based contains matcher cannot have text values",
				con.Key,
			)
		}
		for _, c := range con.Checks.Numeric {
			if c.Min > c.Max {
				probs.AddError(
					src,
					"%s: numeric-based contains matcher has minimum (%f) > maximum (%f)",
					con.Key,
					c.Min,
					c.Max,
				)
			}
		}
	}
}
