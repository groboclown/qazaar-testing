// Under the Apache-2.0 License
package main

import (
	"github.com/groboclown/qazaar-testing/rule-engine/config"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest"
	"github.com/groboclown/qazaar-testing/rule-engine/ingest/sont"
)

type allData struct {
	OntDescriptors *sont.Descriptors
}

func readAll(c *config.ProjectConfig) (*allData, error) {
	ret := allData{
		OntDescriptors: sont.New(),
	}
	ontFiles, err := ingest.FindFiles(c.RefDirs, c.OntologyFiles)
	if err != nil {
		return nil, err
	}
	if err := ingest.ReadOntology(ret.OntDescriptors, ontFiles); err != nil {
		return nil, err
	}

	// TODO documents, rules
	// TODO validate

	return &ret, nil
}
