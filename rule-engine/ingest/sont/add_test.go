// Under the Apache-2.0 License
package sont_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/groboclown/qazaar-testing/rule-engine/ingest/sont"
	"github.com/groboclown/qazaar-testing/rule-engine/schema/ontology"
)

//go:embed "ont-sample.json"
var sample1 []byte

// Test_Add tests adding a new ontology element
//
// This also tests theories about how the json decoding should work.
func Test_Add(t *testing.T) {
	var ont ontology.OntologyV1SchemaJson
	if err := json.Unmarshal(sample1, &ont); err != nil {
		t.Fatal(err)
	}
	if len(ont.Descriptors) != 3 {
		t.Fatalf("Expected 3 descriptors, found %v", ont.Descriptors)
	}

	s := sont.New()
	s.Add(&ont)
	if s.Problems.HasProblems() {
		t.Errorf("Descriptor had problems: %v", s.Problems.Problems())
	}
	if len(s.Enum) != 1 {
		t.Error("Descriptor[0] was not an enum")
	}
	if len(s.Free) != 1 {
		t.Error("Descriptor[1] was not a free")
	}
	if len(s.Numeric) != 1 {
		t.Error("Descriptor[2] was not a numeric")
	}
}
