// Under the Apache-2.0 License
package sont

type TypedDescriptor struct {
	Enum    *EnumDesc
	Free    *FreeDesc
	Numeric *NumericDesc
}

func (s *AllowedDescriptors) Find(key string) *TypedDescriptor {
	t, ok := s.keyTypes[key]
	if !ok {
		return nil
	}
	var e *EnumDesc = nil
	var f *FreeDesc = nil
	var n *NumericDesc = nil
	switch t {
	case EnumDescriptorType:
		e = s.Enum[key]
	case FreeDescriptorType:
		f = s.Free[key]
	case NumericDescriptorType:
		n = s.Numeric[key]
	}
	return &TypedDescriptor{
		Enum:    e,
		Free:    f,
		Numeric: n,
	}
}
