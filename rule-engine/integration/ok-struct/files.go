// Under the Apache-2.0 License
package okstruct_test

import (
	"embed"
	"os"
	"path"

	"github.com/groboclown/qazaar-testing/rule-engine/config"
)

//go:embed "*.json"
var dataFiles embed.FS

func writeFiles(outdir string) error {
	entries, err := dataFiles.ReadDir(".")
	if err != nil {
		return err
	}
	for _, f := range entries {
		if !f.IsDir() {
			data, err := dataFiles.ReadFile(f.Name())
			if err != nil {
				return err
			}
			err = os.WriteFile(path.Join(outdir, f.Name()), data, 0400)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func newConfig(outdir string) *config.ProjectConfig {
	return &config.ProjectConfig{
		LevelMap:      map[string]int{"debug": 0, "info": 1, "warn": 2, "error": 3},
		InfoLevel:     1,
		WarningLevel:  2,
		ErrorLevel:    3,
		RefDirs:       []string{outdir},
		RuleFiles:     []string{"*.rule.json"},
		OntologyFiles: []string{"*.ont.json"},
	}
}
