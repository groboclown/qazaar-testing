// Under the Apache-2.0 License
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/groboclown/qazaar-testing/rule-engine/config"
	"github.com/groboclown/qazaar-testing/rule-engine/engine/runner"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
	"github.com/groboclown/qazaar-testing/rule-engine/validate"
)

var (
	configFile string
	reportDir  string
)

func init() {
	flag.StringVar(&configFile, "config-file", "", "Configuration file location")
	flag.StringVar(&reportDir, "report-dir", "", "Generated report directory")
}

func main() {
	// This ships general error messages to stderr,
	// and the informational problems to stdout (what the end user cares about).
	flag.Parse()

	if configFile == "" {
		fmt.Fprintln(os.Stderr, "Error: must set 'config-file' value.")
		os.Exit(1)
	}
	pc, err := config.ReadProjectConfigFile(configFile)
	if err != nil || pc == nil {
		fmt.Printf("Error reading config file '%s': %s", configFile, err.Error())
		os.Exit(1)
	}

	ctx := context.Background()
	data, validationProbs := ReadValidate(pc, flag.Args(), ctx)
	ReportProblems(validationProbs, os.Stdout)
	if validationProbs.HasErrors() {
		fmt.Fprintf(os.Stderr, "Loading data encountered unrecoverable problems.")
		os.Exit(1)
	}

	engineProbs := RunEngine(pc, data, ctx)
	ReportProblems(engineProbs, os.Stdout)
	if engineProbs.HasErrors() {
		fmt.Fprintf(os.Stderr, "Documents have rule conformity issues.")
		os.Exit(1)
	}

	// Note: this should also create report files in the report directory.
}

func ReadValidate(
	cfg *config.ProjectConfig,
	docFiles []string,
	ctx context.Context,
) (*ingest.AllData, *problem.ProblemSet) {
	probGen, probRead := problem.Async(ctx)

	data := ingest.ReadAll(cfg, docFiles, probGen, ctx)
	validate.ValidateAllDataAsync(data, probGen, ctx)

	probGen.Complete()
	probs := probRead.Read(ctx)
	return data, probs
}

func RunEngine(
	cfg *config.ProjectConfig,
	data *ingest.AllData,
	ctx context.Context,
) *problem.ProblemSet {
	engine := runner.New(data, cfg)
	state, pReader := engine.Start(ctx)
	for state.Step() {
	}
	state.Stop()
	return pReader.Read(ctx)
}
