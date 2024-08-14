// Under the Apache-2.0 License
package okstruct_test

import (
	"context"
	"testing"

	"github.com/groboclown/qazaar-testing/rule-engine/engine/runner"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
	"github.com/groboclown/qazaar-testing/rule-engine/validate"
)

func Test_Engine(t *testing.T) {
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

	pAdder, pReader = problem.Async(ctx)
	validate.ValidateAllDataAsync(all, pAdder, ctx)
	pAdder.Complete()
	probs = pReader.Read(ctx)
	if probs.HasProblems() {
		t.Fatal(probs.Problems())
	}

	engine := runner.New(all, c)
	state, pReader := engine.Start(ctx)
	for state.Step() {
	}
	state.Stop()
	probs = pReader.Read(ctx)
	if probs.HasProblems() {
		t.Fatal(probs.Problems())
	}
}
