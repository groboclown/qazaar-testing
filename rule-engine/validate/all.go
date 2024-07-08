// Under the Apache-2.0 License
package validate

import (
	"context"

	"github.com/groboclown/qazaar-testing/rule-engine/ingest"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
)

func ValidateAllDataAsync(all *ingest.AllData, probs problem.Adder, ctx context.Context) {
	if all == nil {
		return
	}
	chO := ValidateOntologyAsync(all.OntDescriptors, probs, ctx)
	chD := ValidateDocumentsAsync(all.Documents, all.OntDescriptors, probs, ctx)
	chR := ValidateRuleSetAsync(all.RuleSets, all.OntDescriptors, probs, ctx)

	// Completion order doesn't matter.
	<-chO
	<-chD
	<-chR
}
