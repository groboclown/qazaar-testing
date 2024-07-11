// Under the Apache-2.0 License
package obj

func (o *EngineObj) Value(key string) DescriptorValue {
	if e, ok := o.Enum[key]; ok {
		return DescriptorValue{
			Text:   e.List(),
			Number: nil,
		}
	}
	if f, ok := o.Free[key]; ok {
		return DescriptorValue{
			Text:   f.List(),
			Number: nil,
		}
	}
	if n, ok := o.Numeric[key]; ok {
		return DescriptorValue{
			Text:   nil,
			Number: n.List(),
		}
	}
	return DescriptorValue{}
}

func (o EngineObj) Count(key string) int {
	if e, ok := o.Enum[key]; ok {
		return e.Count()
	}
	if f, ok := o.Free[key]; ok {
		return f.Count()
	}
	if n, ok := o.Numeric[key]; ok {
		return n.Count()
	}
	return 0
}
