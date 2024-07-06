// Under the Apache-2.0 License
package descriptor

import "strconv"

type DescriptorValue struct {
	Text   *string
	Number *float64
}

func (d DescriptorValue) String() string {
	if d.Text != nil {
		return *d.Text
	}
	if d.Number != nil {
		return strconv.FormatFloat(*d.Number, 'f', -1, 64)
	}
	return ""
}

type Descriptor struct {
	Key    string
	Text   []string
	Number []float64
}
