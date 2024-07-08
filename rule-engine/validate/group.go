// Under the Apache-2.0 License
package validate

import (
	"context"
	"sync"

	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/descriptor"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/sont"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/srule"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
)

func ValidateGroupAsync(
	group *srule.Group,
	ont *sont.AllowedDescriptors,
	probs problem.Adder,
	ctx context.Context,
) <-chan bool {
	ret := make(chan bool)

	go func() {
		defer func() {
			ret <- onDefer("group", nil, probs)
			close(ret)
		}()

		var wg sync.WaitGroup

		if group != nil {
			ValidateMatchersAsync(group.Matchers, ont, &wg, group.Sources, probs)
			for _, a := range group.Alterations {
				wg.Add(1)
				go func() {
					defer onDefer("group alteration", &wg, probs)
					ValidateAlteration(&a, ont, probs)
				}()
			}
			for _, c := range group.Convergences {
				wg.Add(1)
				go func() {
					defer onDefer("group convergence", &wg, probs)
					ValidateConvergence(&c, ont, probs)
				}()
			}
		}

		wg.Wait()
	}()

	return ret
}

func ValidateAlteration(
	alt *srule.Alteration,
	ont *sont.AllowedDescriptors,
	probs problem.Adder,
) {
	if alt == nil || ont == nil {
		return
	}
	ValidateDescriptor(
		"alteration",
		&descriptor.Descriptor{
			Key:    alt.Key,
			Text:   alt.TextValues,
			Number: alt.NumberValues,
		},
		ont,
		alt.Sources,
		probs,
	)
}

func ValidateConvergence(
	con *srule.Convergence,
	ont *sont.AllowedDescriptors,
	probs problem.Adder,
) {
	if con == nil || ont == nil {
		return
	}
	// Value type of the key not specified, and not needed.
	// However, this can check for the existence of the key.
	checkKey("convergence", con.Key, ont, con.Sources, probs)
}
