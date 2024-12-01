package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	strategyFlag  = flag.String("strategy", "default", "Merging strategy to use.")
	overwriteFlag = flag.Bool("overwrite", false, "Overwrite duplicate keys.")
	skipFlag      = flag.Bool("skip", false, "Skip duplicate keys.")
	outputFlag    = flag.String("output", "", "Output file to write the merged YAML (default: stdout).")
	verboseFlag   = flag.Bool("verbose", false, "Enable more verbose output.")
	skipKeysFlag  = keyList{}
)

func init() {
	flag.Var(&skipKeysFlag, "skipKey", "Key to skip during merging (can be used multiple times).")
}

func main() {
	flag.Usage = printUsage
	flag.Parse()

	files := flag.Args()
	if len(files) < 2 {
		printError("At least two YAML files are required to merge.")
		os.Exit(1)
	}

	// Load and parse YAML files
	maps, err := loadFiles(files)
	if err != nil {
		printError(err.Error())
		os.Exit(2)
	}

	var strategy MergeStrategy
	switch *strategyFlag {
	case "default":
		if *verboseFlag {
			fmt.Println("Using default merging strategy.")
		}
		strategy = &DefaultMergeStrategy{
			Overwrite: *overwriteFlag,
			Skip:      *skipFlag,
			SkipKeys:  skipKeysFlag.toMap(),
		}
	default:
		printError(fmt.Sprintf("Unknown strategy: %s", *strategyFlag))
		os.Exit(3)
	}

	// merge the maps
	merged, err := mergeWithStrategy(strategy, maps...)
	if err != nil {
		printError(fmt.Sprintf("Error during merging: %s", err))
		os.Exit(4)
	}

	// Write the output
	if err := writeOutput(merged, *outputFlag); err != nil {
		printError(fmt.Sprintf("Error writing output: %s", err))
		os.Exit(4)
	}
}

func printUsage() {
	fmt.Println(`======================================================
         YAMYAM (YAM - Yet Another Merger)
       Seamlessly merge multiple YAML files 
======================================================

Usage:
  yamyam [options] <file1.yaml> <file2.yaml> ... 

Options:
  -overwrite       Overwrite duplicate keys during merging.
  -skip            Skip duplicate keys instead of throwing errors.
  -skipKey <key>   Skip specific keys during merging. Can be used multiple times.
                   Example: -skipKey "apiKey" -skipKey "password"
  -output <file>   Specify an output file for the merged YAML.
                   Default: stdout.
  -verbose         Enable more verbose output.

Examples:
  1. Merge files and overwrite duplicates:
     yaml-merge -overwrite file1.yaml file2.yaml

  2. Merge files, skip duplicates, and write to a file:
     yaml-merge -skip -output merged.yaml file1.yaml file2.yaml

  3. Skip specific keys while merging:
     yaml-merge -skipKey "apiKey" -skipKey "debug" file1.yaml file2.yaml

Notes:
  - At least two input files are required for merging.
  - The tool supports merging nested maps and arrays.

======================================================`)
}

// MergeStrategy - Merging strategy interface
type MergeStrategy interface {
	Merge(target, source map[string]any) error
}

// mergeWithStrategy merges multiple maps using the provided strategy
func mergeWithStrategy(strategy MergeStrategy, objs ...map[string]any) (map[string]any, error) {
	result := make(map[string]interface{})
	for _, obj := range objs {
		if err := strategy.Merge(result, obj); err != nil {
			return nil, err
		}
	}
	return result, nil
}
