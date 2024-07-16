// Under the Apache-2.0 License
package concurrent_test

import (
	"context"
	"testing"

	"github.com/groboclown/qazaar-testing/rule-engine/engine/concurrent"
)

func Test_RunEarlyExit(t *testing.T) {
	t.Run("find-deep", func(t *testing.T) {
		ctx := context.Background()
		tree := R1{
			v: "a",
			r: []R1{
				{v: "b", r: []R1{{v: "b1"}}},
				{v: "c"},
			},
		}
		res := <-concurrent.RunEarlyExit(rEarlyConst, ctx, []R1{tree})
		if res == nil || *res != "b1" {
			t.Errorf("did not return 'b1', but %v", res)
		}
		if err := ctx.Err(); err != nil {
			t.Error(err)
		}
	})
	t.Run("find-not-exist", func(t *testing.T) {
		ctx := context.Background()
		r := UnorderedSearch{10}
		values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		res := <-concurrent.RunEarlyExit(r, ctx, [][]int{values})
		if res != nil {
			t.Errorf("incorrectly marked a returned value %v", res)
		}
		if err := ctx.Err(); err != nil {
			t.Error(err)
		}
	})
	t.Run("find-exist", func(t *testing.T) {
		ctx := context.Background()
		r := UnorderedSearch{6}
		values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		res := <-concurrent.RunEarlyExit(r, ctx, [][]int{values})
		if res == nil || *res != true {
			t.Errorf("did not return 'true', but %v", res)
		}
		if err := ctx.Err(); err != nil {
			t.Error(err)
		}
	})
}

func Benchmark_RunsWell_Deep(b *testing.B) {
	ctx := context.Background()
	tree := R1{
		v: "a",
		r: []R1{
			{v: "b", r: []R1{{v: "b1"}}},
			{v: "c"},
		},
	}

	b.StartTimer()
	res := <-concurrent.RunEarlyExit(rEarlyConst, ctx, []R1{tree})
	b.StopTimer()

	if res == nil || *res != "b1" {
		b.Errorf("did not return 'b1', but %v", res)
	}
	if err := ctx.Err(); err != nil {
		b.Error(err)
	}
}

func Benchmark_RunsWell_Found(b *testing.B) {
	r := UnorderedSearch{6}
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	ctx := context.Background()
	b.StartTimer()
	res := <-concurrent.RunEarlyExit(r, ctx, [][]int{values})
	b.StopTimer()
	if res == nil || *res != true {
		b.Errorf("did not return 'true', but %v", res)
	}
	if err := ctx.Err(); err != nil {
		b.Error(err)
	}
}

func Benchmark_RunsWell_NotFound(b *testing.B) {
	r := UnorderedSearch{10}
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	ctx := context.Background()
	b.StartTimer()
	res := <-concurrent.RunEarlyExit(r, ctx, [][]int{values})
	b.StopTimer()
	if res != nil {
		b.Errorf("incorrectly returned non-nil value %v", res)
	}
	if err := ctx.Err(); err != nil {
		b.Error(err)
	}
}

type R1 struct {
	v string
	r []R1
}

type REarly struct{}

var rEarlyConst REarly

func (r REarly) Perform(v R1, o chan<- R1) (*string, bool) {
	if len(v.v) > 1 {
		return &v.v, true
	}
	for _, x := range v.r {
		o <- x
	}
	return nil, false
}

type UnorderedSearch struct{ search int }

var RetTrue = true

func (b UnorderedSearch) Perform(v []int, o chan<- []int) (*bool, bool) {
	vl := len(v)
	if vl <= 0 {
		return nil, false
	}
	var mid int = len(v) / 2
	if v[mid] == b.search {
		return &RetTrue, true
	}
	if mid > 0 {
		o <- v[0 : mid-1]
	}
	if mid+1 < vl {
		o <- v[mid+1:]
	}
	return nil, false
}
