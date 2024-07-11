// Under the Apache-2.0 License
package descriptor

import "github.com/groboclown/qazaar-testing/rule-engine/schema/rules"

func DecodeRuleValues(vals []rules.AlterationValuesElem) []DescriptorValue {
	ret := make([]DescriptorValue, len(vals))
	for i, v := range vals {
		ret[i] = Decode(v)
	}
	return ret
}
