// Under the Apache-2.0 License
package descriptor

import "strings"

type DescriptorValueTypes interface{ string | float64 }

type DescriptorValueTypeTransform[T DescriptorValueTypes] interface {
	Transform(val T) T
}

// DistinctMap allows for constructing values that obey the 'distinct' requirement.
func DistinctMap[T DescriptorValueTypes]() map[T]bool {
	return make(map[T]bool)
}

// DistinctMapArray turns the distinct map into an array.
func DistinctMapArray[T DescriptorValueTypes](m map[T]bool) []T {
	ret := make([]T, 0, len(m))
	for k := range m {
		ret = append(ret, k)
	}
	return ret
}

// DescriptorValueCopy creates a copy of the array argument.
func DescriptorValueCopy[T DescriptorValueTypes](a []T, transform *DescriptorValueTypeTransform[T]) []T {
	ret := make([]T, len(a))
	if transform == nil {
		copy(ret, a)
	} else {
		for i, v := range a {
			ret[i] = (*transform).Transform(v)
		}
	}
	return ret
}

// DistinctValueArray turns a list of values (numeric or text) into an array of just distinct values.
func DistinctValueArray[T DescriptorValueTypes](in []T, transform DescriptorValueTypeTransform[T]) []T {
	m := make(map[T]bool)
	for _, v := range in {
		m[transform.Transform(v)] = true
	}
	ret := make([]T, len(m))
	i := 0
	for k := range m {
		ret[i] = k
		i++
	}
	return ret
}

// AddToDistinctMap adds a list of values to a map of values.
//
// The `transform` argument allows for handling case sensitivity requirements.
func AddToDistinctMap[T DescriptorValueTypes](
	base map[T]bool, adders []T, transform DescriptorValueTypeTransform[T],
) {
	for _, k := range adders {
		base[transform.Transform(k)] = true
	}
}

// AppendToList extends the base list of values with the adder list.
//
// The `transform` argument allows for handling case sensitivity requirements.
// This assumes the base values have already passed through the transform.
func AppendToList[T DescriptorValueTypes](
	base []T, adders []T, transform *DescriptorValueTypeTransform[T],
) []T {
	if transform == nil {
		if base == nil {
			if adders == nil {
				return make([]T, 0)
			}
			return adders
		}
		return append(base, adders...)
	}
	basel := len(base)
	if basel == 0 {
		if adders == nil {
			return make([]T, 0)
		}
	}
	ret := make([]T, basel+len(adders))
	copy(ret, base)
	for i, v := range adders {
		ret[basel+i] = (*transform).Transform(v)
	}
	return ret
}

type floatDT int

func (f floatDT) Transform(val float64) float64 {
	return val
}

type strDT int

func (s strDT) Transform(val string) string {
	return val
}

type strLowerDT int

func (s strLowerDT) Transform(val string) string {
	return strings.ToLower(val)
}

var FloatTransform DescriptorValueTypeTransform[float64] = floatDT(0)
var FloatTransformP *DescriptorValueTypeTransform[float64] = nil
var StringTransform DescriptorValueTypeTransform[string] = strDT(1)
var StringTransformP *DescriptorValueTypeTransform[string] = nil
var StringLowerTransform DescriptorValueTypeTransform[string] = strLowerDT(2)
var StringLowerTransformP *DescriptorValueTypeTransform[string] = &StringLowerTransform
