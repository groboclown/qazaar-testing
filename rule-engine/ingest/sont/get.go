// Under the Apache-2.0 License
package sont

type TypedDescriptor struct {
	Enum    *EnumDesc
	Free    *FreeDesc
	Numeric *NumericDesc

	key     string
	keyType DescriptorType
}

func (t *TypedDescriptor) Type() DescriptorType {
	return t.keyType
}
func (t *TypedDescriptor) KeyName() string {
	return t.key
}
func (t *TypedDescriptor) IsDistinct() bool {
	switch {
	case t.Enum != nil:
		return t.Enum.Distinct
	case t.Free != nil:
		return t.Free.Distinct
	case t.Numeric != nil:
		return t.Numeric.Distinct
	default:
		return false
	}
}
func (t *TypedDescriptor) IsCaseSensitive() bool {
	if t.Free != nil {
		return t.Free.CaseSensitive
	}
	// Otherwise, default to an exact match.
	return true
}

func (s *AllowedDescriptors) Find(key string) *TypedDescriptor {
	t, ok := s.keyTypes[key]
	if !ok || t == UnknownDescriptorType {
		return nil
	}
	var e *EnumDesc = nil
	var f *FreeDesc = nil
	var n *NumericDesc = nil
	switch t {
	case EnumDescriptorType:
		e = s.enums[key]
	case FreeDescriptorType:
		f = s.frees[key]
	case NumericDescriptorType:
		n = s.numerics[key]
	}
	return &TypedDescriptor{
		key:     key,
		keyType: t,
		Enum:    e,
		Free:    f,
		Numeric: n,
	}
}

func (s *AllowedDescriptors) Type(key string) DescriptorType {
	t, ok := s.keyTypes[key]
	if !ok {
		return UnknownDescriptorType
	}
	return t
}

func (s *AllowedDescriptors) Enums() []*EnumDesc {
	return valueArray(s.enums)
}

func (s *AllowedDescriptors) Frees() []*FreeDesc {
	return valueArray(s.frees)
}

func (s *AllowedDescriptors) Numerics() []*NumericDesc {
	return valueArray(s.numerics)
}

func valueArray[T EnumDesc | FreeDesc | NumericDesc](m map[string]*T) []*T {
	ret := make([]*T, 0, len(m))
	for _, v := range m {
		ret = append(ret, v)
	}
	return ret
}

func (s *AllowedDescriptors) Enum(key string) *EnumDesc {
	v, ok := s.enums[key]
	if !ok {
		return nil
	}
	return v
}

func (s *AllowedDescriptors) Free(key string) *FreeDesc {
	v, ok := s.frees[key]
	if !ok {
		return nil
	}
	return v
}

func (s *AllowedDescriptors) Numeric(key string) *NumericDesc {
	v, ok := s.numerics[key]
	if !ok {
		return nil
	}
	return v
}

func (d *EnumDesc) Type() DescriptorType {
	return EnumDescriptorType
}

func (d *EnumDesc) KeyName() string {
	return d.Key
}

func (d *EnumDesc) IsDistinct() bool {
	return d.Distinct
}

func (d *FreeDesc) Type() DescriptorType {
	return EnumDescriptorType
}

func (d *FreeDesc) KeyName() string {
	return d.Key
}

func (d *FreeDesc) IsDistinct() bool {
	return d.Distinct
}

func (d *NumericDesc) Type() DescriptorType {
	return EnumDescriptorType
}

func (d *NumericDesc) KeyName() string {
	return d.Key
}

func (d *NumericDesc) IsDistinct() bool {
	return d.Distinct
}
