// Under the Apache-2.0 License
package matcher_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/groboclown/qazaar-testing/rule-engine/engine/matcher"
	"github.com/groboclown/qazaar-testing/rule-engine/engine/obj"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/sont"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/srule"
)

func Test_IsContainsMatch(t *testing.T) {
	ont, err := ingest.ParseOntology(strings.NewReader(`
	{
		"$schema": "",
		"descriptors": [
			{
				"type":         "number",
				"key":          "e1",
				"distinct":     false,
				"maximum":      100,
				"minimum":      -100,
				"maximumCount": 100
			}
		]
	}
	`), "test")
	if err != nil {
		t.Fatal(err)
	}
	descriptors := sont.New()
	descriptors.Add(ont)
	if descriptors.Problems.HasProblems() {
		t.Fatal(descriptors.Problems.Problems())
	}
	factory := obj.NewObjFactory(descriptors)
	source := obj.ObjSource{}

	t.Run("1-number-all-ok", func(t *testing.T) {
		ob := factory.Empty(source)
		ob.Add("e1", floats(1.5))
		c := &srule.ContainsMatcher{
			Operation: srule.ContainsAll,
			Count:     false,
			Distinct:  false,
			Key:       "e1",
			Checks: srule.ValueCheckSet{
				Numeric: []srule.NumericBoundsCheck{{Min: 1, Max: 2}},
				// Text: []srule.StringCheck{{R: regexp.MustCompile("^a+$")}},
			},
		}
		o := ob.Seal()
		if ok, errs := matcher.IsContainsMatch(o, c); !ok {
			t.Errorf("did not match: %v", errs)
		}
	})
	t.Run("1-number-all-fail", func(t *testing.T) {
		ob := factory.Empty(source)
		ob.Add("e1", floats(1.5))
		c := &srule.ContainsMatcher{
			Operation: srule.ContainsAll,
			Count:     false,
			Distinct:  false,
			Key:       "e1",
			Checks: srule.ValueCheckSet{
				Numeric: []srule.NumericBoundsCheck{{Min: 2, Max: 3}},
				// Text: []srule.StringCheck{{R: regexp.MustCompile("^a+$")}},
			},
		}
		o := ob.Seal()
		ok, errs := matcher.IsContainsMatch(o, c)
		if ok {
			t.Fatal("incorrectly matched")
		}
		if diff := diffMatchers(errs, []matcher.MatcherMismatch{
			{Obj: o, Contains: &matcher.MismatchContains{c}},
		}); diff != "" {
			t.Error(diff)
		}
	})
	t.Run("2.2-number-all-ok", func(t *testing.T) {
		ob := factory.Empty(source)
		ob.Add("e1", floats(1.5, 1.8))
		c := &srule.ContainsMatcher{
			Operation: srule.ContainsAll,
			Count:     false,
			Distinct:  false,
			Key:       "e1",
			Checks: srule.ValueCheckSet{
				Numeric: []srule.NumericBoundsCheck{{Min: 1, Max: 2}, {Min: 1, Max: 2}},
				// Text: []srule.StringCheck{{R: regexp.MustCompile("^a+$")}},
			},
		}
		o := ob.Seal()
		if ok, errs := matcher.IsContainsMatch(o, c); !ok {
			t.Errorf("did not match: %v", errs)
		}
	})
	t.Run("2.2-number-all-fail", func(t *testing.T) {
		ob := factory.Empty(source)
		ob.Add("e1", floats(2.2, -1.0))
		c := &srule.ContainsMatcher{
			Operation: srule.ContainsAll,
			Count:     false,
			Distinct:  false,
			Key:       "e1",
			Checks: srule.ValueCheckSet{
				Numeric: []srule.NumericBoundsCheck{{Min: 1, Max: 2}, {Min: 1, Max: 2}},
				// Text: []srule.StringCheck{{R: regexp.MustCompile("^a+$")}},
			},
		}
		o := ob.Seal()
		ok, errs := matcher.IsContainsMatch(o, c)
		if ok {
			t.Fatal("incorrectly matched")
		}
		if diff := diffMatchers(errs, []matcher.MatcherMismatch{
			{Obj: o, Contains: &matcher.MismatchContains{c}},
		}); diff != "" {
			t.Error(diff)
		}
	})
	t.Run("2.2-number-all-partial-fail", func(t *testing.T) {
		ob := factory.Empty(source)
		ob.Add("e1", floats(2.2, 1.1))
		c := &srule.ContainsMatcher{
			Operation: srule.ContainsAll,
			Count:     false,
			Distinct:  false,
			Key:       "e1",
			Checks: srule.ValueCheckSet{
				Numeric: []srule.NumericBoundsCheck{{Min: 1, Max: 2}, {Min: 1, Max: 2}},
				// Text: []srule.StringCheck{{R: regexp.MustCompile("^a+$")}},
			},
		}
		o := ob.Seal()
		ok, err := matcher.IsContainsMatch(o, c)
		if ok {
			t.Fatal("incorrectly matched")
		}
		if diff := diffMatchers(err, []matcher.MatcherMismatch{
			{Obj: o, Contains: &matcher.MismatchContains{c}},
		}); diff != "" {
			t.Error(diff)
		}
	})
	t.Run("2.1-number-exact-partial-fail", func(t *testing.T) {
		ob := factory.Empty(source)
		ob.Add("e1", floats(2.2, 1.1))
		c := &srule.ContainsMatcher{
			Operation: srule.ContainsExactly,
			Count:     false,
			Distinct:  false,
			Key:       "e1",
			Checks: srule.ValueCheckSet{
				Numeric: []srule.NumericBoundsCheck{{Min: 1, Max: 2}, {Min: 1, Max: 2}},
				// Text: []srule.StringCheck{{R: regexp.MustCompile("^a+$")}},
			},
		}
		o := ob.Seal()

		ok, err := matcher.IsContainsMatch(o, c)
		if ok {
			t.Fatal("incorrectly matched")
		}
		if diff := diffMatchers(err, []matcher.MatcherMismatch{
			{Obj: o, Collection: nil, Contains: &matcher.MismatchContains{c}},
		}); diff != "" {
			t.Error(diff)
		}
	})
}

func floats(v ...float64) obj.DescriptorValues {
	return obj.DescriptorValues{Number: v}
}

func diffMatchers(actual, expected []matcher.MatcherMismatch) string {
	actStr := make([]string, len(actual))
	for i, a := range actual {
		actStr[i] = a.String()
	}
	expStr := make([]string, len(expected))
	for i, e := range expected {
		expStr[i] = e.String()
	}
	return cmp.Diff(actStr, expStr)
}
