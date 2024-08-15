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
		_, e := fmt.Fprint(out, "Quiet:\n")
		errs = append(errs, e)
		for _, p := range quiet {
			_, e := fmt.Fprintf(out, "  %s\n", p.String())
			errs = append(errs, e)
		}
	}
	if len(info) > 0 {
		_, e := fmt.Fprint(out, "Informative:\n")
		errs = append(errs, e)
		for _, p := range info {
			_, e := fmt.Fprintf(out, "  %s\n", p.String())
			errs = append(errs, e)
		}
	}
	if len(warn) > 0 {
		_, e := fmt.Fprint(out, "Warnings:\n")
		errs = append(errs, e)
		for _, p := range warn {
			_, e := fmt.Fprintf(out, "  %s\n", p.String())
			errs = append(errs, e)
		}
	}
	if len(err) > 0 {
		_, e := fmt.Fprint(out, "Errors:\n")
		errs = append(errs, e)
		for _, p := range err {
			_, e := fmt.Fprintf(out, "  %s\n", p.String())
			errs = append(errs, e)
		}
	}

	return errors.Join(errs...)
}
