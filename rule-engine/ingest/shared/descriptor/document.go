package descriptor

import "github.com/groboclown/qazaar-testing/rule-engine/schema/document"

func DecodeDocumentValues(vals []document.DocumentDescriptorValuesElem) []DescriptorValue {
	ret := make([]DescriptorValue, len(vals))
	for i, v := range vals {
		ret[i] = Decode(v)
	}
	return ret
}

func JoinDocumentDescriptors(descs []document.DocumentDescriptor) []*Descriptor {
	ret := make([]*Descriptor, len(descs))
	for i, d := range descs {
		t, n := Join(DecodeDocumentValues(d.Values))
		ret[i] = JoinKeyValues(string(d.Key), t, n)
	}
	return ret
}
