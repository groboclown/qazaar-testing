// Under the Apache-2.0 License
package descriptor_test

import (
	"testing"

	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/descriptor"
)

func Test_DistinctValueArray(t *testing.T) {
	t.Run("string-no-dups", diffDescTest(
		descriptor.JoinKeyValues("k1", []string{"a", "b", "c"}, nil),
		"k1",
		[]string{"a", "b", "c"},
		nil,
	))
	t.Run("number-no-dups", diffDescTest(
		descriptor.JoinKeyValues("k2", nil, []float64{1, 2.2, 2.222, 3}),
		"k2",
		nil,
		[]float64{1, 2.2, 2.222, 3},
	))
	t.Run("string-dups", diffDescTest(
		descriptor.JoinKeyValues("k3", []string{"aa", "aa", "bb", "bb"}, nil),
		"k3",
		[]string{"aa", "bb"},
		nil,
	))
	t.Run("number-dups", diffDescTest(
		descriptor.JoinKeyValues("k3", nil, []float64{0, 1.1111, 1.1111, 2, -2e100, -2e100}),
		"k3",
		nil,
		[]float64{0, 1.1111, 2, -2e100},
	))
}

func diffDescTest(actual *descriptor.Descriptor, key string, text []string, num []float64) func(t *testing.T) {
	return func(t *testing.T) {
		diffDesc(t, actual, key, text, num)
	}
}

func diffDesc(t *testing.T, actual *descriptor.Descriptor, key string, text []string, num []float64) {
	// Order of items doesn't matter.
	if actual == nil {
		t.Fatal("returned nil descriptor")
	}
	if actual.Key != key {
		t.Errorf("expected key '%s', found '%s'", key, actual.Key)
	}
	diffSet(t, "text", actual.Text, text)
	diffSet(t, "numeric", actual.Number, num)
}

func diffSet[T float64 | string](t *testing.T, name string, act []T, exp []T) {
	actSet := asSet(act)
	expSet := asSet(exp)
	actExp := asSet(act)
	for k := range expSet {
		if _, ok := actSet[k]; !ok {
			t.Errorf("%s did not include %v", name, k)
		} else {
			delete(actExp, k)
		}
	}
	for k := range actExp {
		t.Errorf("%s included extra item %v", name, k)
	}
}

func asSet[T float64 | string](in []T) map[T]bool {
	ret := make(map[T]bool)
	for _, k := range in {
		ret[k] = true
	}
	return ret
}
