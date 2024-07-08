// Under the Apache-2.0 License
package sdoc_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/groboclown/qazaar-testing/rule-engine/ingest/sdoc"
	"github.com/groboclown/qazaar-testing/rule-engine/schema/document"
)

//go:embed "doc-sample.json"
var sample1 []byte

func Test_Add(t *testing.T) {
	var doc document.DocumentDescriptionV1SchemaJson
	if err := json.Unmarshal(sample1, &doc); err != nil {
		t.Fatal(err)
	}
	if len(doc.Objects) != 1 {
		t.Fatalf("expected 1 document object, found %v", doc)
	}

	s := sdoc.New()
	s.Add(&doc)
	if s.Problems.HasProblems() {
		t.Fatalf("Loading the document encountered problems: %v", s.Problems.Problems())
	}
	if len(s.Objects) != 1 {
		t.Fatalf("expected 1 added object, found %v", s.Objects)
	}

	obj := s.Objects[0]
	if obj == nil {
		t.Fatalf("added object is nil")
	}
	if len(obj.Descriptors) != 2 {
		t.Errorf("expected 2 descriptors, found %v", obj.Descriptors)
	}
	if obj.Id != "d1" {
		t.Errorf("expected id 'd1', found '%s'", obj.Id)
	}
	if len(obj.Sources) != 1 {
		t.Fatalf("expected 1 source, found %v", obj.Sources)
	}

	if obj.Descriptors[0].Key != "req1" {
		t.Errorf("descriptor[0] key should = 'req1', found '%s'", obj.Descriptors[0].Key)
	}
	if len(obj.Descriptors[0].Number) != 0 {
		t.Errorf("descriptor[0] should have only text, found numbers %v", obj.Descriptors[0].Number)
	}
	if len(obj.Descriptors[0].Text) != 2 || obj.Descriptors[0].Text[0] != "a" || obj.Descriptors[0].Text[1] != "b" {
		t.Errorf("descriptor[0] text should be ['a','b'], found %v", obj.Descriptors[0].Text)
	}

	if obj.Descriptors[1].Key != "req2" {
		t.Errorf("descriptor[1] key should = 'req1', found '%s'", obj.Descriptors[1].Key)
	}
	if len(obj.Descriptors[1].Text) != 0 {
		t.Errorf("descriptor[1] should have only numbers, found text %v", obj.Descriptors[0].Text)
	}
	if len(obj.Descriptors[1].Number) != 2 || obj.Descriptors[1].Number[0] != 1.1 || obj.Descriptors[1].Number[1] != 2 {
		t.Errorf("descriptor[1] number should be [1.1, 2], found %v", obj.Descriptors[0].Text)
	}
}
