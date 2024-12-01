package main

import (
	"fmt"
	"github.com/goccy/go-yaml"
	"os"
	"strings"
)

// Print an error message
func printError(message string) {
	_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", message)
}

// Loads and parses YAML files
func loadFiles(files []string) ([]map[string]any, error) {
	var maps []map[string]any
	for _, file := range files {
		if *verboseFlag {
			fmt.Printf("Loading file: %s\n", file)
		}
		m, err := parseFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to parse file '%s': %w", file, err)
		}
		maps = append(maps, m)
	}
	return maps, nil
}

// Writes the merged map to the specified output
func writeOutput(data map[string]any, outputPath string) error {
	writer := os.Stdout
	if outputPath != "" {
		f, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)
		writer = f
	}
	return yaml.NewEncoder(writer).Encode(data)
}

// Parses a YAML file into a map
func parseFile(path string) (map[string]any, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	var data map[string]any
	if err := yaml.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

// keyList holds a list of keys to skip
type keyList []string

func (kl *keyList) String() string {
	return strings.Join(*kl, ",")
}

func (kl *keyList) Set(value string) error {
	*kl = append(*kl, value)
	return nil
}

// Converts the key list to a map for faster lookup
func (kl *keyList) toMap() map[string]struct{} {
	skipMap := make(map[string]struct{}, len(*kl))
	for _, key := range *kl {
		skipMap[key] = struct{}{}
	}
	return skipMap
}
