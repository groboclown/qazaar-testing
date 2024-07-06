// Under the Apache-2.0 License
package comments

import "github.com/groboclown/qazaar-testing/rule-engine/schema/document"

func JoinDocComments(obj *document.DocumentObject) []string {
	ret := make([]string, 0)
	if obj != nil {
		if obj.Comment != nil && *obj.Comment != "" {
			ret = append(ret, string(*obj.Comment))
		}
		for _, c := range obj.Comments {
			if c != "" {
				ret = append(ret, string(c))
			}
		}
	}
	return ret
}
