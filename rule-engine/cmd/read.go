// Under the Apache-2.0 License
package main

import (
	"context"

	"github.com/groboclown/qazaar-testing/rule-engine/config"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/sdoc"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/sont"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/srule"
	"github.com/groboclown/qazaar-testing/rule-engine/problem"
	"github.com/groboclown/qazaar-testing/rule-engine/schema/document"
	"github.com/groboclown/qazaar-testing/rule-engine/schema/ontology"
	"github.com/groboclown/qazaar-testing/rule-engine/schema/rules"
)

type allData struct {
	OntDescriptors *sont.AllowedDescriptors
	RuleSets       *srule.RuleSet
	Documents      *sdoc.Documents
}

// readAll reads and parses all the data asynchronously.
func readAll(
	c *config.ProjectConfig,
	docFiles []string,
	probs problem.Adder,
	ctx context.Context,
) *allData {

	ret := allData{
		OntDescriptors: sont.New(),
		RuleSets:       srule.New(),
		Documents:      sdoc.New(),
	}

	ont := readOnt(c, probs, ctx)
	rule := readRule(c, probs, ctx)
	doc := readDocument(docFiles, probs, ctx)

	ontDone := false
	ruleDone := false
	docDone := false
	for {
		select {
		case o, ok := <-ont:
			if !ok {
				ontDone = true
			}
			ret.OntDescriptors.Add(o)
		case r, ok := <-rule:
			if !ok {
				ruleDone = true
			}
			ret.RuleSets.Add(r)
		case d, ok := <-doc:
			if !ok {
				docDone = true
			}
			ret.Documents.Add(d)
		case <-ctx.Done():
			ontDone = true
			ruleDone = true
			docDone = true
		}

		if ontDone && ruleDone && docDone {
			break
		}
	}

	return &ret
}

func readOnt(
	c *config.ProjectConfig,
	probs problem.Adder,
	ctx context.Context,
) <-chan *ontology.OntologyV1SchemaJson {
	ret := make(chan *ontology.OntologyV1SchemaJson)

	go func() {
		defer close(ret)
		ch := ingest.FindFilesAsync(c.RefDirs, c.OntologyFiles, ctx)
		for {
			select {
			case f, ok := <-ch:
				if !ok {
					return
				}
				ont, err := ingest.ReadOntologyFile(f)
				if err != nil {
					probs.Error(f, err)
				}
				if ont != nil {
					ret <- ont
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return ret
}

func readRule(
	c *config.ProjectConfig,
	probs problem.Adder,
	ctx context.Context,
) <-chan *rules.RulesV1SchemaJson {
	ret := make(chan *rules.RulesV1SchemaJson)

	go func() {
		defer close(ret)
		ch := ingest.FindFilesAsync(c.RefDirs, c.RuleFiles, ctx)
		for {
			select {
			case f, ok := <-ch:
				if !ok {
					return
				}
				rule, err := ingest.ReadRuleFile(f)
				if err != nil {
					probs.Error(f, err)
				}
				if rule != nil {
					ret <- rule
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return ret
}

func readDocument(
	files []string,
	probs problem.Adder,
	ctx context.Context,
) <-chan *document.DocumentDescriptionV1SchemaJson {
	ret := make(chan *document.DocumentDescriptionV1SchemaJson)

	go func() {
		defer close(ret)

		for _, f := range files {
			if ctx.Err() != nil {
				return
			}
			doc, err := ingest.ReadDocumentsFile(f)
			if err != nil {
				probs.Error(f, err)
			}
			if doc != nil {
				ret <- doc
			}
		}
	}()

	return ret
}
