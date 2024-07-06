// Under the Apache-2.0 License
package sources

func Join(s1 []Source, s2 ...Source) []Source {
	la := len(s1)
	ret := make([]Source, la+len(s2))
	_ = copy(ret, s1)
	_ = copy(ret[la:], s2)
	return ret
}
