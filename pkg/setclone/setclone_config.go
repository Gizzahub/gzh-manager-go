package setclone

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path"
)

type SetcloneGithub struct {
	RootPath string `yaml:"root_path"`
	Provider string `yaml:"provider"`
	url      string `yaml:"url"`
	Protocol string `yaml:"protocol"`
	OrgName  string `yaml:"org_name"`
}

type SetcloneGitlab struct {
	RootPath  string `yaml:"root_path"`
	Provider  string `yaml:"provider"`
	Url       string `yaml:"url"`
	Recursive bool   `yaml:"recursive"`
	GroupName string `yaml:"group_name"`
	Protocol  string `yaml:"protocol"`
}

type SetcloneDefault struct {
	Protocol string         `yaml:"protocol"`
	Github   SetcloneGithub `yaml:"github"`
	Gitlab   SetcloneGitlab `yaml:"gitlab"`
}

type setcloneConfig struct {
	Version           string           `yaml:"version"`
	Default           SetcloneDefault  `yaml:"default"`
	IgnoreNameRegexes []string         `yaml:"ignore_names"`
	RepoRoots         []SetcloneGithub `yaml:"repo_roots"`
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
