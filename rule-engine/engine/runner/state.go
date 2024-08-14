// Under the Apache-2.0 License
package runner

import (
	"sync"

	"github.com/groboclown/qazaar-testing/rule-engine/engine/obj"
	"github.com/groboclown/qazaar-testing/rule-engine/engine/sog"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/srule"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
)

type engineRunnerState struct {
	engine   *engineRunner
	sogs     []*sog.SogBuilder
	problems problem.Adder
	objects  []*obj.EngineObj
}

func (s *engineRunnerState) Stop() {
	s.problems.Complete()
}

func (s *engineRunnerState) Step() bool {
	// The whole thing is a SOG build & gather tool.
	newObjCh := make(chan *obj.EngineObj)
	joinedCh := joinObjAsync(newObjCh)

	var wg sync.WaitGroup

	for _, builder := range s.sogs {
		wg.Add(1)
		go func(builder *sog.SogBuilder) {
			defer func() {
				wg.Done()
				s.problems.Recover("engineRunner.Step", recover())
			}()

			// First step - match the objects against the SOGs.
			// FIXME MAJOR BUG INFINITE LOOPING.
			// The Add should indeed pass over every object, but it should
			// only allow the builder to add the object if the builder added
			// based on a new object.  If an old object creates a new sog, then
			// it should not be allowed.
			builder.Reset()
			for _, o := range s.objects {
				builder.Add(o)
			}

			// Once the SOGs are gathered...
			for _, si := range builder.Seal() {
				// Match members against the Convergence.
				for _, c := range si.Group().Convergences {
					wg.Add(1)
					go func(m []*obj.EngineObj, c *srule.Convergence) {
						defer func() {
							wg.Done()
							s.problems.Recover("engineRunner.Step.Convergence", recover())
						}()
						if v := MatchConvergence(si.Group().Id, m, c, s.engine.ont); v != nil {
							s.problems.Add(convAsProblem(s.engine.levelMap, v))
						}
					}(si.Members(), &c)
				}

				// Match the new SOG values against the rules.
				o := si.Obj()

				// Add the SOG values into the objects.
				newObjCh <- o
			}
		}(builder)
	}

	// Wait for the async build to complete, then wrap up the async.
	wg.Wait()
	close(newObjCh)
	addedObj := <-joinedCh
	s.objects = append(s.objects, addedObj...)

	// Return 'false' if no more SOG objects were created.
	added := len(addedObj) > 0
	if !added {
		s.Stop()
	}
	return added
}

func joinObjAsync(c <-chan *obj.EngineObj) <-chan []*obj.EngineObj {
	done := make(chan []*obj.EngineObj)
	objs := make([]*obj.EngineObj, 0)
	go func() {
		defer func() {
			done <- objs
			close(done)
		}()
		for o := range c {
			objs = append(objs, o)
		}
	}()
	return done
}
