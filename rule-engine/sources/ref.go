// Under the Apache-2.0 License
package sources

import "fmt"

func universalRefId(loc, rep string, ver *string) string {
	v := ""
	if ver != nil {
		v = *ver
	}
	return fmt.Sprintf("%s|%s|%s", loc, rep, v)
}

func (sg *SourceGen) addRef(loc, rep string, ver *string) *innerRef {
	if sg == nil {
		return nil
	}
	id := universalRefId(loc, rep, ver)
	if r, ok := sg.refs[id]; ok && r != nil {
		return r
	}
	r := &innerRef{
		loc: loc,
		rep: rep,
		ver: ver,
	}
	sg.refs[id] = r
	return r
}
