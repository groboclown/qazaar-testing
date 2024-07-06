// Under the Apache-2.0 License
package comments

import "github.com/groboclown/qazaar-testing/rule-engine/schema/ontology"

func JoinOntComments(com *ontology.Comment, cl ontology.CommentList) []string {
	ret := make([]string, 0)
	if com != nil && *com != "" {
		ret = append(ret, string(*com))
	}
	for _, c := range cl {
		if c != "" {
			ret = append(ret, string(c))
		}
	}
	return ret
}
