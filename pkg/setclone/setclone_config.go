package setclone

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path"
)

type setcloneConfig struct {
	GithubUrl         string   `yaml:"github_url"`
	GithubToken       string   `yaml:"github_token"`
	DefaultProtocol   string   `yaml:"default_protocol"`
	IgnoreNameRegexes []string `yaml:"ignore_names"`
}

// Assuming these helper functions are defined within the same package
func getConfigDir() string {
	// Implement this function to return the configuration directory path
	return "/path/to/config"
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func (cfg *setcloneConfig) ConfigExists(targetPath string) bool {
	return fileExists(path.Join(targetPath, "setclone.yaml"))
}

func (cfg *setcloneConfig) ReadConfig(targetPath string) {
	configPath := path.Join(targetPath, "setclone.yaml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		log.Fatalf("failed to unmarshal config file: %v", err)
	}
}
