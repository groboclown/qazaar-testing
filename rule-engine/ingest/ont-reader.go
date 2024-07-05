// Under the Apache-2.0 License
package ingest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/groboclown/qazaar-testing/rule-engine/ingest/sont"
	"github.com/groboclown/qazaar-testing/rule-engine/schema/ontology"
)

// ReadOntology adds all the ontology files listed in the project configuration.
func ReadOntology(d *sont.Descriptors, files []string) error {
	errs := make([]error, 0)
	for _, f := range files {
		src, err := ReadOntologyFile(f)
		if err != nil {
			errs = append(errs, err)
		}
		d.Add(src)
	}
	return errors.Join(errs...)
}

func ReadOntologyFile(f string) (*ontology.OntologyV1SchemaJson, error) {
	r, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	return ParseOntology(r, f)
}

func ParseOntology(r io.Reader, src string) (*ontology.OntologyV1SchemaJson, error) {
	var ret ontology.OntologyV1SchemaJson
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", src, err.Error())
	}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", src, err.Error())
	}
	return &ret, nil
}
