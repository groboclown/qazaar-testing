// Under the Apache-2.0 License
package okstruct_test

import (
	"context"
	"testing"

	"github.com/groboclown/qazaar-testing/rule-engine/ingest"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
)

func Test_Reader(t *testing.T) {
	tmp := t.TempDir()
	if err := writeFiles(tmp); err != nil {
		t.Fatal(err)
	}
	c := newConfig(tmp)
	ctx := context.Background()
	probAdd, reader := problem.Async(ctx)
	all := ingest.ReadAll(c, []string{}, probAdd, ctx)
	probAdd.Complete()

	probs := reader.Read(ctx)
	if probs.HasProblems() {
		t.Fatal(probs.Problems())
	}

	if len(all.Documents.Objects) != 0 {
		t.Errorf("incorrectly read documents")
	}
}
