// Under the Apache-2.0 License
package config

import "github.com/groboclown/qazaar-testing/rule-engine/problem"

// ProjectConfig defines a project setup for processing the rules.
type ProjectConfig struct {
	LevelMap      map[string]int `json:"level-map"` // Maps between level names and a error level; 0 is lowest
	InfoLevel     int            `json:"info"`      // Level at or above for informative issues.
	WarningLevel  int            `json:"warn"`      // Level at or above for warnings.
	ErrorLevel    int            `json:"error"`     // Level at or above for errors.
	RefDirs       []string       `json:"ref-dir"`   // Base directory for finding the rule and ontology files
	RuleFiles     []string       `json:"rules"`     // Glob pattern for rule files under the ref dirs
	OntologyFiles []string       `json:"ontology"`  // Glob pattern for ontology files under the ref dirs
}

// RuntimeConfig contains shared data for processing the rules.
type RuntimeConfig struct {
	Problems *problem.ProblemSet
}
