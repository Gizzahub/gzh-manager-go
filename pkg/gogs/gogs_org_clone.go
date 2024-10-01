package gogs

import (
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type setcloneConfig struct {
	GithubUrl       string   `yaml:"github_url"`
	GithubToken     string   `yaml:"github_token"`
	DefaultProtocol string   `yaml:"default_protocol"`
	IgnoreNames     []string `yaml:"ignore_names"`
}

// getConfigDir returns the configuration directory path
func getConfigDir() string {
	// Implement this function to return the configuration directory path
	return "/path/to/config"
}

// fileExists checks if a file exists at the given path
func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// ConfigExists checks if the configuration file exists
func (cfg *setcloneConfig) ConfigExists() bool {
	configDir := getConfigDir()
	return fileExists(path.Join(configDir, "setclone.yaml"))
}

// ReadConfig reads and unmarshals the YAML configuration file
func (cfg *setcloneConfig) ReadConfig(fullpath string) {
	configDir := getConfigDir()
	configPath := path.Join(configDir, fullpath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		log.Fatalf("failed to unmarshal config file: %v", err)
	}
}
