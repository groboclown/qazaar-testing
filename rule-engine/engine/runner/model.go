// Under the Apache-2.0 License
package runner

import (
	"context"

	"github.com/groboclown/qazaar-testing/rule-engine/problem"
)

// EngineRunner allows high-level access to the engine processing system.
type EngineRunner interface {
	Start(ctx context.Context) (EngineState, problem.ProblemConsumer)
}

type EngineState interface {
	// Step returns 'false' if it encounters an end state.
	Step() bool
	Stop()
}
