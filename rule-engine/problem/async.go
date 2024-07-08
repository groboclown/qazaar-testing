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
	complete *ProblemSet
	done     <-chan *ProblemSet
}

func Async(ctx context.Context) (Adder, ProblemConsumer) {
	ch := make(chan *Problem)
	done := make(chan *ProblemSet)

	go func() {
		defer close(done)
		ret := New()

		for {
			select {
			case p, ok := <-ch:
				if p != nil {
					ret.Add(*p)
				}
				if !ok {
					done <- ret
					return
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
				close(ch)
				done <- ret
				return
			}
		}
	}()

	// TODO start the consumer in a go function, so it doesn't block the generator.
	return &asyncProblemGenerator{out: ch}, &asyncProblemConsumer{done: done}
}

func (pg *asyncProblemGenerator) Complete() {
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

func (pg *asyncProblemGenerator) Recover(source string, recover any) {
	if recover != nil {
		pg.Add(Problem{
			Sources: nil,
			Level:   Err,
			Message: fmt.Sprintf("%s: runtime error (%v)", source, recover),
		})
	}
}

func (pg *asyncProblemGenerator) Add(p ...Problem) {
	if pg == nil {
		return
	}
	for _, v := range p {
		pg.out <- &v
	}
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
	if c == nil {
		return nil
	}
	if c.complete != nil {
		return c.complete
	}

	select {
	case p := <-c.done:
		c.complete = p
		return p
	case <-ctx.Done():
		return nil
	}
}
