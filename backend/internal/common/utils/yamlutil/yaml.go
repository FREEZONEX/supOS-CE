package yamlutil

import (
	"fmt"
	"io"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// LoadYamlFromFile reads a YAML file from the given path and unmarshals it into an object of the specified generic type P.
// It returns a pointer to the newly created object, which provides a more type-safe and convenient usage.
func LoadYamlFromFile[P any](path string) (*P, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read yaml file %s: %w", path, err)
	}

	var out P
	err = yaml.Unmarshal(data, &out)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml from %s: %w", path, err)
	}

	return &out, nil
}

// LoadYamlFromReader reads YAML data from an io.Reader and unmarshals it into an object of the specified generic type P.
// This is the equivalent of the Java version that accepts an InputStream.
func LoadYamlFromReader[P any](reader io.Reader) (*P, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read from reader: %w", err)
	}

	var out P
	err = yaml.Unmarshal(data, &out)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml from reader: %w", err)
	}

	return &out, nil
}

// WriteTempYamlFile marshals the provided model object into YAML and writes it to a temporary file.
// The temporary file will have a name pattern like "fileName_1678886400.yml".
// It returns the created file. The caller is responsible for cleaning up the file.
func WriteTempYamlFile[P any](fileName string, model P) (*os.File, error) {
	// Marshal the object to YAML bytes
	data, err := yaml.Marshal(model)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal model to yaml: %w", err)
	}

	// Create a temporary file
	pattern := fmt.Sprintf("%s_%d_*.yml", fileName, time.Now().Unix())
	file, err := os.CreateTemp("", pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	// Write the YAML data to the file
	if _, err := file.Write(data); err != nil {
		// Attempt to close and remove the file on error
		file.Close()
		os.Remove(file.Name())
		return nil, fmt.Errorf("failed to write yaml to temp file: %w", err)
	}

	return file, nil
}
