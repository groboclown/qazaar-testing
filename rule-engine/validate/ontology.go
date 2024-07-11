// Under the Apache-2.0 License
package validate

import (
	"context"
	"sync"

	"github.com/groboclown/qazaar-testing/rule-engine/ingest/sont"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
)

// ValidateOntologyAsync checks the ontology against known restrictions.
func ValidateOntologyAsync(
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

		if ont == nil {
			return
		}

		var wg sync.WaitGroup

		for _, d := range ont.Enums() {
			wg.Add(1)
			go func() {
				defer onDefer("ontology enum", &wg, probs)
				ValidateOntEnum(d, probs)
			}()
		}
		for _, d := range ont.Frees() {
			wg.Add(1)
			go func() {
				defer onDefer("ontology free", &wg, probs)
				ValidateOntFree(d, probs)
			}()
		}
		for _, d := range ont.Numerics() {
			wg.Add(1)
			go func() {
				defer onDefer("ontology numeric", &wg, probs)
				ValidateOntNumeric(d, probs)
			}()
		}

		wg.Wait()
	}()

	return ret
}

func ValidateOntEnum(d *sont.EnumDesc, probs problem.Adder) {
	// Nothing really to do here.
}

func ValidateOntFree(d *sont.FreeDesc, probs problem.Adder) {
	// Nothing really to do here.
}

func ValidateOntNumeric(d *sont.NumericDesc, probs problem.Adder) {
	if d == nil {
		return
	}
	if d.Minimum > d.Maximum {
		probs.AddError(
			d.Sources,
			"%s: numeric descriptor has minimum (%f) > maximum (%f)",
			d.Key,
			d.Minimum,
			d.Maximum,
		)
	}
	if d.Distinct {
		// Need to research if this is really the case.
		probs.AddInfo(
			d.Sources,
			"%s: numeric descriptor defined as distinct; be careful with this - for non-integral values, this can lead to unexpected behavior.",
			d.Key,
		)
	}
}
