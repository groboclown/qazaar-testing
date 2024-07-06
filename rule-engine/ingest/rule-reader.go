// Under the Apache-2.0 License
package ingest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/groboclown/qazaar-testing/rule-engine/ingest/srule"
	"github.com/groboclown/qazaar-testing/rule-engine/schema/rules"
)

// ReadRule adds all the rule files listed in the project configuration.
func ReadRule(d *srule.RuleSet, files []string) error {
	errs := make([]error, 0)
	for _, f := range files {
		src, err := ReadRuleFile(f)
		if err != nil {
			errs = append(errs, err)
		}
		d.Add(src)
	}
	return errors.Join(errs...)
}

func ReadRuleFile(f string) (*rules.RulesV1SchemaJson, error) {
	r, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	return ParseRule(r, f)
}

func ParseRule(r io.Reader, src string) (*rules.RulesV1SchemaJson, error) {
	var ret rules.RulesV1SchemaJson
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
