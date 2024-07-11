// Under the Apache-2.0 License
package okstruct_test

import (
	"context"
	"testing"

	"github.com/groboclown/qazaar-testing/rule-engine/ingest"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
	"github.com/groboclown/qazaar-testing/rule-engine/validate"
)

func Test_Reader(t *testing.T) {
	tmp := t.TempDir()
	if err := writeFiles(tmp); err != nil {
		t.Fatal(err)
	}
	c := newConfig(tmp)
	docs := []string{
		mustDocFilename(tmp, "om-source", t),
		mustDocFilename(tmp, "openapi-source", t),
		mustDocFilename(tmp, "sql-source", t),
	}
	ctx := context.Background()
	pAdder, pReader := problem.Async(ctx)
	all := ingest.ReadAll(c, docs, pAdder, ctx)

	pAdder.Complete()
	probs := pReader.Read(ctx)
	probs.Add(all.Problems().Problems()...)
	if probs.HasProblems() {
		t.Fatal(probs.Problems())
	}

	if len(all.OntDescriptors.Enums()) != 3 {
		t.Errorf("incorrectly read ont-enum (%d)", len(all.OntDescriptors.Enums()))
	}
	if len(all.OntDescriptors.Numerics()) != 0 {
		t.Errorf("incorrectly read ont-numeric (%d)", len(all.OntDescriptors.Numerics()))
	}
	if len(all.OntDescriptors.Frees()) != 4 {
		t.Errorf("incorrectly read ont-free (%d)", len(all.OntDescriptors.Frees()))
	}
	if len(all.RuleSets.Groups) != 2 {
		t.Errorf("incorrectly read rules-groups (%d)", len(all.RuleSets.Groups))
	}
	if len(all.RuleSets.Rules) != 2 {
		t.Errorf("incorrectly read rules-rules (%d)", len(all.RuleSets.Rules))
	}
	if len(all.Documents.Objects) != 13 {
		t.Errorf("incorrectly read documents (%d)", len(all.Documents.Objects))
	}

	pAdder, pReader = problem.Async(ctx)
	validate.ValidateAllDataAsync(all, pAdder, ctx)

	pAdder.Complete()
	probs = pReader.Read(ctx)
	if probs.HasProblems() {
		t.Fatal(probs.Problems())
	}
}
