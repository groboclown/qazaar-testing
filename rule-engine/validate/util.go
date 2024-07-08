// Under the Apache-2.0 License
//
// Basic helper tools.
package validate

import (
	"sync"

	"github.com/groboclown/qazaar-testing/rule-engine/problem"
)

func onDefer(context string, wg *sync.WaitGroup, probs problem.Adder) bool {
	r := recover()
	probs.Recover(context, r)
	if wg != nil {
		wg.Done()
	}
	return r == nil
}
