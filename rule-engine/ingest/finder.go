// Under the Apache-2.0 License
package ingest

import (
	"context"
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

// FindFilesAsync scans all the root directories for all matching files.
func FindFilesAsync(roots []string, globs []string, ctx context.Context) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		done := make(chan int)

		for i, r := range roots {
			sub := FindRootedFilesAsync(r, globs, ctx)
			go func() {
				defer func() { done <- i }()

				for {
					select {
					case f, ok := <-sub:
						if !ok {
							return
						}
						ch <- f
					case <-ctx.Done():
						return
					}
				}
			}()
		}

		for range roots {
			select {
			case <-done:
			case <-ctx.Done():
				close(done)
				return
			}
		}
	}()
	return ch
}

// FindRootedFilesAsync scans the root directory for all matching files.
func FindRootedFilesAsync(root string, globs []string, ctx context.Context) <-chan string {
	_, cancel := context.WithCancelCause(ctx)
	ch := make(chan string)
	go func() {
		defer close(ch)

		for _, g := range globs {
			for strings.HasPrefix(g, "/") {
				g = strings.TrimPrefix(g, "/")
			}
			if g != "" {
				p := path.Join(root, g)
				m, err := filepath.Glob(p)
				if err != nil {
					cancel(err)
					return
				}
				for _, f := range m {
					ch <- f
				}
			}
		}
	}()

	return ch
}
