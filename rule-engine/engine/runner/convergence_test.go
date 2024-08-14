// Test the convergence implementation.
//
// Under the Apache-2.0 License
package runner

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/groboclown/qazaar-testing/rule-engine/engine/obj"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/descriptor"
)

func Test_loadMaps(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		nm := make(map[float64]int)
		tm := make(map[string]int)

		loadMaps(tm, nm, obj.DescriptorValues{}, descriptor.StringTransform)

		if len(nm) != 0 {
			t.Errorf("Loaded numbers: %v", nm)
		}
		if len(tm) != 0 {
			t.Errorf("Loaded text: %v", tm)
		}
	})

	t.Run("singles", func(t *testing.T) {
		nm := make(map[float64]int)
		tm := make(map[string]int)

		loadMaps(tm, nm, obj.DescriptorValues{
			Text:   []string{"a", "b"},
			Number: []float64{1, 2.2},
		}, descriptor.StringTransform)

		if diff := cmp.Diff(
			map[string]int{"a": 1, "b": 1},
			tm,
		); diff != "" {
			t.Errorf("bad text: %s", diff)
		}

		if diff := cmp.Diff(
			map[float64]int{1: 1, 2.2: 1},
			nm,
		); diff != "" {
			t.Errorf("bad text: %s", diff)
		}
	})

	t.Run("many", func(t *testing.T) {
		nm := make(map[float64]int)
		tm := make(map[string]int)

		loadMaps(tm, nm, obj.DescriptorValues{
			Text:   []string{"a", "a", "a", "b", "b", "a"},
			Number: []float64{1, 2.2, 1, 1, 2.2, 2.3},
		}, descriptor.StringTransform)

		if diff := cmp.Diff(
			map[string]int{"a": 4, "b": 2},
			tm,
		); diff != "" {
			t.Errorf("bad text: %s", diff)
		}

		if diff := cmp.Diff(
			map[float64]int{1: 3, 2.2: 2, 2.3: 1},
			nm,
		); diff != "" {
			t.Errorf("bad text: %s", diff)
		}
	})
}
