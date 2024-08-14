// Under the Apache-2.0 License
package runner

import (
	"context"

	"github.com/groboclown/qazaar-testing/rule-engine/config"
	"github.com/groboclown/qazaar-testing/rule-engine/engine/obj"
	"github.com/groboclown/qazaar-testing/rule-engine/engine/sog"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/sont"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/srule"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
)

// New turns the documents into engine objects.
func New(data *ingest.AllData, config *config.ProjectConfig) EngineRunner {
	if data == nil {
		return nil
	}
	factory := obj.NewObjFactory(data.OntDescriptors)
	base := make([]*obj.EngineObj, len(data.Documents.Objects))
	for i, o := range data.Documents.Objects {
		base[i] = factory.FromDocument(o)
	}
	return &engineRunner{
		factory:  factory,
		ont:      data.OntDescriptors,
		groups:   data.RuleSets.Groups,
		rules:    data.RuleSets.Rules,
		base:     base,
		levelMap: convertLevelMap(config),
	}
}

type engineRunner struct {
	factory  obj.ObjFactory                  // Shared factory.
	ont      *sont.AllowedDescriptors        // Ontology.
	groups   []*srule.Group                  // All Self-Organizing Group definitions.
	rules    []*srule.Rule                   // All rules.
	base     []*obj.EngineObj                // Initial set of document objects.
	levelMap map[string]problem.ProblemLevel // Maps the rule level to a problem level.
}

func (e *engineRunner) Start(ctx context.Context) (EngineState, problem.ProblemConsumer) {
	// Load all the initial rule issues against the base objects.  That way,
	// each step will only check the assembled SOGs for violations, and there won't
	// be duplicate rule checks.
	adder, consumer := problem.Async(ctx)
	addRuleProblems(adder, e.levelMap, checkAllAgainstRules(e.base, e.rules))

	sogs := make([]*sog.SogBuilder, len(e.groups))
	for i, g := range e.groups {
		sogs[i] = sog.NewBuilder(g, e.factory)
	}

	return &engineRunnerState{
		engine:   e,
		objects:  e.base,
		sogs:     sogs,
		problems: adder,
	}, consumer
}

func convertLevelMap(config *config.ProjectConfig) map[string]problem.ProblemLevel {
	ret := make(map[string]problem.ProblemLevel)
	for key, level := range config.LevelMap {
		switch {
		case level < config.InfoLevel:
			ret[key] = problem.Quiet
		case level < config.WarningLevel:
			ret[key] = problem.Info
		case level <= config.ErrorLevel:
			ret[key] = problem.Warn
		default:
			ret[key] = problem.Err
		}
	}
	return ret
}
