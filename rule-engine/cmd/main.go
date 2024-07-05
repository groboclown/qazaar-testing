// Under the Apache-2.0 License
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/groboclown/qazaar-testing/rule-engine/config"
)

var (
	configFile string
	reportDir  string
)

func init() {
	flag.StringVar(&configFile, "config-file", "", "Configuration file location")
	flag.StringVar(&reportDir, "report-dir", "", "Generated report directory")
}

func main() {
	if configFile == "" {
		fmt.Println("Error: must set 'config-file' value.")
		os.Exit(1)
	}
	pc, err := config.ReadProjectConfigFile(configFile)
	if err != nil || pc == nil {
		fmt.Printf("Error reading config file '%s': %s", configFile, err.Error())
		os.Exit(1)
	}

	_, err = readAll(pc)
	if err != nil {
		fmt.Printf("Error reading configuration data: %s", err.Error())
		os.Exit(1)
	}
}
