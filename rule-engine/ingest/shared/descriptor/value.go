// Under the Apache-2.0 License
package descriptor

// Join turns the list of shared descriptor value objects into a split list of string and float.
//
// This does not handle the "distinct key" case.
func Join(vals []DescriptorValue) (sl []string, fl []float64) {
	sl = make([]string, 0)
	fl = make([]float64, 0)
	for _, v := range vals {
		if v.Text != nil {
			sl = append(sl, *v.Text)
		}
		if v.Number != nil {
			fl = append(fl, *v.Number)
		}
	}
	return
}

// JoinKey turns the key and list of values into a shared descriptor.
//
// This does not handle the "distinct key" case.
func JoinKey(key string, vals []DescriptorValue) *Descriptor {
	sl, fl := Join(vals)
	return JoinKeyValues(key, sl, fl)
}

// JoinKeyValues turns the key and list of values into a shared descriptor.
//
// This does not handle the "distinct key" case.  For that, use `DistinctValueArray()` or
// `Descriptor.Distinct()`.
func JoinKeyValues(key string, text []string, numbers []float64) *Descriptor {
	return &Descriptor{
		Key:    key,
		Text:   text,
		Number: numbers,
	}
}

// Decode turns the simple value (string or numeric) into a DescriptorValue.
func Decode(val any) DescriptorValue {
	var f float64 = 0
	switch v := val.(type) {
	case string:
		return DescriptorValue{Text: &v}
	case int:
		f = float64(v)
	case int8:
		f = float64(v)
	case int16:
		f = float64(v)
	case int32:
		f = float64(v)
	case int64:
		f = float64(v)
	case uint8:
		f = float64(v)
	case uint16:
		f = float64(v)
	case uint32:
		f = float64(v)
	case uint64:
		f = float64(v)
	case float32:
		f = float64(v)
	case float64:
		f = v
	}
	return DescriptorValue{Number: &f}
}
