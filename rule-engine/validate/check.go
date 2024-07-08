// Under the Apache-2.0 License
package validate

import (
	"fmt"
	"regexp"

	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/descriptor"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/sources"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/sont"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
	"github.com/groboclown/qazaar-testing/rule-engine/schema/ontology"
)

func ValidateDescriptor(
	context string,
	d *descriptor.Descriptor,
	ont *sont.AllowedDescriptors,
	src []sources.Source,
	probs problem.Adder,
) {
	if d == nil || ont == nil {
		return
	}
	typed := checkKey(context, d.Key, ont, src, probs)
	if typed == nil {
		return
	}
	// These validations also include the source for the offending ontology.
	if typed.Enum != nil {
		ValidateEnum(context, d, typed.Enum, sources.Join(src, typed.Enum.Sources...), probs)
	}
	if typed.Free != nil {
		ValidateFree(context, d, typed.Free, sources.Join(src, typed.Free.Sources...), probs)
	}
	if typed.Numeric != nil {
		ValidateNumeric(context, d, typed.Numeric, sources.Join(src, typed.Numeric.Sources...), probs)
	}
}

func ValidateEnum(
	context string,
	d *descriptor.Descriptor,
	ont *sont.EnumDesc,
	src []sources.Source,
	probs problem.Adder,
) {
	if d == nil || ont == nil {
		return
	}
	if len(d.Number) != 0 {
		probs.AddError(
			src,
			"%s: enum %s cannot have numeric values",
			d.Key,
			context,
		)
	}
	checkCount(context, d.Key, d.Text, ont.MaximumCount, src, probs)
	for _, t := range d.Text {
		if _, ok := ont.Enum[t]; !ok {
			probs.AddError(
				src,
				"%s: enum %s invalid value (%s)",
				d.Key,
				context,
				t,
			)
		}
	}
	if ont.Distinct {
		findDuplicates(context, d.Key, d.Text, src, probs)
	}
}

func ValidateFree(
	context string,
	d *descriptor.Descriptor,
	ont *sont.FreeDesc,
	src []sources.Source,
	probs problem.Adder,
) {
	if d == nil || ont == nil {
		return
	}
	if len(d.Number) != 0 {
		probs.AddError(
			src,
			"%s: free %s cannot have numeric values",
			d.Key,
			context,
		)
	}
	checkCount(context, d.Key, d.Text, ont.MaximumCount, src, probs)
	for _, t := range d.Text {
		if len(t) > ont.MaximumLength {
			probs.AddError(
				src,
				"%s: free %s value length (%d) exceeds maximum (%d) (%s)",
				d.Key,
				context,
				len(t),
				ont.MaximumLength,
				t,
			)
		}
		for _, con := range ont.Constraints {
			checkConstraint(
				d.Key,
				t,
				&con,
				src,
				probs,
			)
		}
	}
	if ont.Distinct {
		findDuplicates(context, d.Key, d.Text, src, probs)
	}
}

func ValidateNumeric(
	context string,
	d *descriptor.Descriptor,
	ont *sont.NumericDesc,
	src []sources.Source,
	probs problem.Adder,
) {
	if d == nil || ont == nil {
		return
	}
	if len(d.Text) != 0 {
		probs.AddError(
			src,
			"%s: numeric %s cannot have text values",
			d.Key,
			context,
		)
	}
	checkCount(context, d.Key, d.Number, ont.MaximumCount, src, probs)
	for _, n := range d.Number {
		if n < ont.Minimum || n > ont.Maximum {
			probs.AddError(
				src,
				"%s: numeric %s value (%f) outside bounds [%f, %f]",
				n,
				context,
				ont.Minimum,
				ont.Maximum,
			)
		}
	}
}

func checkKey(
	owningType string,
	key string,
	ont *sont.AllowedDescriptors,
	src []sources.Source,
	probs problem.Adder,
) *sont.TypedDescriptor {
	typed := ont.Find(key)
	if typed == nil {
		probs.AddError(
			src,
			"undefined %s key (%s)",
			owningType,
			key,
		)
		return nil
	}
	return typed
}

func findDuplicates(
	context string,
	key string,
	values []string,
	src []sources.Source,
	probs problem.Adder,
) {
	discovered := make(map[string]bool)
	for _, v := range values {
		if reported, ok := discovered[v]; ok {
			if !reported {
				probs.AddError(
					src,
					"%s: %s does not allow duplicate values (%s)",
					key,
					context,
					v,
				)
				discovered[v] = true
			}
		} else {
			discovered[v] = false
		}
	}
}

func checkCount[T string | float64](
	context string,
	key string,
	values []T,
	maxCount int,
	src []sources.Source,
	probs problem.Adder,
) {
	if len(values) > maxCount {
		probs.AddError(
			src,
			"%s: %s can have a maximum of %d values (found %d)",
			key,
			context,
			maxCount,
			len(values),
		)
	}
}

func checkConstraint(
	key string,
	val string,
	con *sont.ValueConstraint,
	src []sources.Source,
	probs problem.Adder,
) {
	if con == nil {
		return
	}
	switch con.Type {
	case ontology.ValueConstraintTypeFormat:
		probs.AddWarning(
			src,
			"%s: value constraint type '%s' not supported.",
			key,
			ontology.ValueConstraintTypeFormat,
		)
	case ontology.ValueConstraintTypePattern:
		if con.Pattern == nil {
			probs.AddError(
				src,
				"%s: invalid value constraint; no pattern for 'pattern' type",
				key,
			)
			return
		}
		re, err := regexp.Compile(*con.Pattern)
		if err != nil {
			probs.Error(
				fmt.Sprintf("%s: invalid value constraint pattern '%s'", key, *con.Pattern),
				err,
			)
			return
		}
		if !re.Match([]byte(val)) {
			probs.AddError(
				src,
				"%s: value (%s) does not match constraint pattern (%s)",
				key,
				val,
				*con.Pattern,
			)
		}
	}
}
