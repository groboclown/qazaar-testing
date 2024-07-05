// Under the Apache-2.0 License
package ingest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/groboclown/qazaar-testing/rule-engine/ingest/sdoc"
	"github.com/groboclown/qazaar-testing/rule-engine/schema/document"
)

// ReadDocuments adds all the documents files listed in the runtime configuration to the documents object.
func ReadDocuments(d *sdoc.Documents, files []string) error {
	errs := make([]error, 0)
	for _, f := range files {
		src, err := ReadDocumentsFile(f)
		if err != nil {
			errs = append(errs, err)
		}
		d.Add(src)
	}
	return errors.Join(errs...)
}

func ReadDocumentsFile(f string) (*document.DocumentDescriptionV1SchemaJson, error) {
	r, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	return ParseDocuments(r, f)
}

func ParseDocuments(r io.Reader, src string) (*document.DocumentDescriptionV1SchemaJson, error) {
	var ret document.DocumentDescriptionV1SchemaJson
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
