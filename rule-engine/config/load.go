// Under the Apache-2.0 License
package config

import (
	"encoding/json"
	"io"
	"os"
)

// ReadProjectConfigFile reads the project configuration from the file.
func ReadProjectConfigFile(cname string) (*ProjectConfig, error) {
	r, err := os.Open(cname)
	if err != nil {
		return nil, err
	}
	return ParseProjectConfig(r)
}

// ParseProjectConfig parses the project configuration from the reader.
func ParseProjectConfig(reader io.Reader) (*ProjectConfig, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var ret *ProjectConfig
	if err := json.Unmarshal(data, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}
