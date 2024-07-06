// Under the Apache-2.0 License
package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/groboclown/qazaar-testing/rule-engine/problem"
)

func ReportProblems(probs *problem.ProblemSet, out io.Writer) error {
	if probs == nil {
		return nil
	}
	errs := make([]error, 0)

	quiet := probs.ProblemsAt(problem.Quiet)
	info := probs.ProblemsAt(problem.Info)
	warn := probs.ProblemsAt(problem.Warn)
	err := probs.ProblemsAt(problem.Err)

	_, e := fmt.Fprintf(
		out,
		"Quiet Messages: %d\nInformative Messages: %d\nWarning Messages: %d\nError Messages: %d\n",
		len(quiet), len(info), len(warn), len(err),
	)
	errs = append(errs, e)

	if len(quiet) > 0 {
		fmt.Fprint(out, "Quiet:\n")
		for _, p := range quiet {
			fmt.Fprintf(out, "  %s\n", p.String())
		}
	}
	if len(info) > 0 {
		fmt.Fprint(out, "Informative:\n")
		for _, p := range info {
			fmt.Fprintf(out, "  %s\n", p.String())
		}
	}
	if len(warn) > 0 {
		fmt.Fprint(out, "Warnings:\n")
		for _, p := range warn {
			fmt.Fprintf(out, "  %s\n", p.String())
		}
	}
	if len(err) > 0 {
		fmt.Fprint(out, "Errors:\n")
		for _, p := range err {
			fmt.Fprintf(out, "  %s\n", p.String())
		}
	}

	return errors.Join(errs...)
}
