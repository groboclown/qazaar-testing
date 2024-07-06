// Under the Apache-2.0 License
package problem

import (
	"context"
	"fmt"

	"github.com/groboclown/qazaar-testing/rule-engine/ingest/shared/sources"
)

type asyncProblemGenerator struct {
	out chan<- *Problem
}

type ProblemConsumer interface {
	Read(ctx context.Context) *ProblemSet
}

type asyncProblemConsumer struct {
	in <-chan *Problem
}

func Async() (Adder, ProblemConsumer) {
	ch := make(chan *Problem)
	return &asyncProblemGenerator{out: ch}, &asyncProblemConsumer{in: ch}
}

func (pg *asyncProblemGenerator) Done() {
	if pg != nil {
		close(pg.out)
	}
}

func (pg *asyncProblemGenerator) Error(source string, err ...error) {
	for _, e := range err {
		if e != nil {
			pg.Add(Problem{
				Sources: nil,
				Level:   Err,
				Message: fmt.Sprintf("%s: %s", source, e.Error()),
			})
		}
	}
}

func (pg *asyncProblemGenerator) Add(p Problem) {
	if pg == nil {
		return
	}
	pg.out <- &p
}

func (pg *asyncProblemGenerator) AddError(
	sources []sources.Source,
	format string,
	args ...any,
) {
	pg.AddProblem(sources, Err, format, args...)
}

func (pg *asyncProblemGenerator) AddWarning(
	sources []sources.Source,
	format string,
	args ...any,
) {
	pg.AddProblem(sources, Warn, format, args...)
}

func (pg *asyncProblemGenerator) AddInfo(
	sources []sources.Source,
	format string,
	args ...any,
) {
	pg.AddProblem(sources, Info, format, args...)
}

func (pg *asyncProblemGenerator) AddProblem(
	sources []sources.Source,
	level ProblemLevel,
	format string,
	args ...any,
) {
	pg.Add(Problem{
		Level:   level,
		Message: fmt.Sprintf(format, args...),
		Sources: sources,
	})
}

func (c *asyncProblemConsumer) Read(ctx context.Context) *ProblemSet {
	ret := New()
	for {
		select {
		case p, ok := <-c.in:
			if p != nil {
				ret.Add(*p)
			}
			if !ok {
				return ret
			}
		case <-ctx.Done():
			err := context.Cause(ctx)
			if err != nil {
				ret.AddError(
					nil,
					"internal error: %s",
					err.Error(),
				)
			}
			// Should close the channel...
			return ret
		}
	}
}
