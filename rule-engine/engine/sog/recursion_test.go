// Under the Apache-2.0 License
package sog_test

import (
	"testing"

	"github.com/groboclown/qazaar-testing/rule-engine/engine/obj"
	"github.com/groboclown/qazaar-testing/rule-engine/engine/sog"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/sources"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/srule"
)

func Test_IsRecursion(t *testing.T) {
	group := srule.Group{
		Id:        "1",
		Variables: make(map[string]*srule.VariableDef),
		Matchers: &srule.MatchingDescriptorSet{
			// Empty matcher, so it matches anything into the same group.
			Collection: []srule.CollectionMatcher{},
			Contains:   []srule.ContainsMatcher{},
		},
		KeySharedValues: []string{},
		Alterations:     []srule.Alteration{},
		Convergences:    []srule.Convergence{},
		Comments:        []string{},
		Sources:         []sources.Source{},
	}
	s0 := mkSrc()
	s1 := mkSrc()
	s0a := mkSrc(&s0)
	s01ab := mkSrc(&s0a, &s1)
	obj0 := obj.EngineObj{
		Source: s0,
	}
	obj0a := obj.EngineObj{
		Source: s0a,
	}
	obj01ab := obj.EngineObj{
		Source: s01ab,
	}

	// Create a new item from a group
	builder := sog.NewBuilder(&group, obj.NewObjFactory(nil))
	// ... Create the group
	if res := builder.Add(&obj0); res != sog.Created {
		t.Errorf("Expected created result, found %d", res)
	}
	// ... Add another item to the group
	if res := builder.Add(&obj0a); res != sog.Added {
		t.Errorf("Expected created result, found %d", res)
	}

	builtObj := builder.Seal()
	if len(builtObj) != 1 {
		t.Fatalf("expected 1 item from the group, found %v", builtObj)
	}

	t.Run("not-recursive", func(t *testing.T) {
		builder := sog.NewBuilder(&group, obj.NewObjFactory(nil))
		// Add the previous result.
		if res := builder.Add(builtObj[0].Obj()); res != sog.Created {
			t.Errorf("Expected non-recursive result, found %d", res)
		}
		// Add a non-recursive item to the group, but with shared history.
		if res := builder.Add(&obj01ab); res != sog.Added {
			t.Errorf("Expected non-recursive result, found %d", res)
		}

		// Ensure the result contains all added items...
	})
	t.Run("yes-recursive", func(t *testing.T) {
		builder := sog.NewBuilder(&group, obj.NewObjFactory(nil))
		// Add the previous result.
		if res := builder.Add(builtObj[0].Obj()); res != sog.Created {
			t.Errorf("Expected non-recursive result, found %d", res)
		}
		// Add a recursive item to the group.
		if res := builder.Add(&obj0); res != sog.Recursion {
			t.Errorf("Expected recursive result, found %d", res)
		}

		// Ensure the result does not contain the recursive item...
	})
}

func mkSrc(parents ...*obj.ObjSource) obj.ObjSource {
	return obj.ObjSource{
		Source:  []sources.Source{},
		Parents: parents,
	}
}
