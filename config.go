package aime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	// collection of system prompts where key is the systemPromptKey and value is the systemPrompt
	System map[string]string `json:"system"`
}

func getExtension(filename string) string {
	i := strings.LastIndex(filename, ".")
	if i == -1 {
		return ""
	}
	return filename[i+1:]
}

func LoadConfig() (*Config, error) {
	var home string
	if h, err := os.UserHomeDir(); err == nil {
		home = h
	}
	var configFiles = []string{
		".aimerc.json",
		".aimerc.yaml",
		".aimerc.yml",
	}
	if home != "" {
		configFiles = append(
			configFiles,
			filepath.Join(home, ".aimerc.json"),
			filepath.Join(home, ".aimerc.yaml"),
			filepath.Join(home, ".aimerc.yml"),
		)
	}

	var filename string

	for _, f := range configFiles {
		if _, err := os.Stat(f); !os.IsNotExist(err) {
			filename = f
			break
		}
	}
	if filename == "" {
		return nil, nil
	}
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	ext := getExtension(filename)
	var config = new(Config)
	switch ext {
	case "json":
		if err := json.Unmarshal(b, config); err != nil {
			return nil, err
		}
	case "yml", "yaml":
		if err := yaml.Unmarshal(b, config); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unexpected config file extension: %s", ext)
	}

	return config, nil
}
