// Under the Apache-2.0 License
package ingest

import (
	"path"
	"path/filepath"
	"strings"
)

// FindFiles scans all the root directories for all matching files.
func FindFiles(roots []string, globs []string) ([]string, error) {
	ret := make([]string, 0)
	for _, r := range roots {
		m, err := FindRootedFiles(r, globs)
		if err != nil {
			return ret, err
		}
		ret = append(ret, m...)
	}
	return ret, nil
}

// FindRootedFiles scans the root directory for all matching files.
func FindRootedFiles(root string, globs []string) ([]string, error) {
	ret := make([]string, 0)
	for _, g := range globs {
		for strings.HasPrefix(g, "/") {
			g = strings.TrimPrefix(g, "/")
		}
		if g != "" {
			p := path.Join(root, g)
			m, err := filepath.Glob(p)
			if err != nil {
				return ret, err
			}
			ret = append(ret, m...)
		}
	}
	return ret, nil
}
