package setclone

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

// >>>>>>>>>> default >>>>>>>>>>
type SetcloneDefault struct {
	Protocol string                `yaml:"protocol" validate:"required,oneof=http https ssh"`
	Github   SetcloneDefaultGithub `yaml:"github"`
	Gitlab   SetcloneDefaultGitlab `yaml:"gitlab"`
}

type SetcloneDefaultGithub struct {
	RootPath string `yaml:"root_path"`
	Provider string `yaml:"provider"`
	url      string `yaml:"url"`
	Protocol string `yaml:"protocol"`
	OrgName  string `yaml:"org_name"`
}

type SetcloneDefaultGitlab struct {
	RootPath  string `yaml:"root_path"`
	Provider  string `yaml:"provider"`
	Url       string `yaml:"url"`
	Recursive bool   `yaml:"recursive"`
	Protocol  string `yaml:"protocol"`
	GroupName string `yaml:"group_name"`
}

// <<<<<<<<<< default <<<<<<<<<<

type SetcloneGithub struct {
	RootPath string `yaml:"root_path" validate:"required"`
	Provider string `yaml:"provider" validate:"required"`
	url      string `yaml:"url"`
	Protocol string `yaml:"protocol" validate:"required,oneof=http https ssh"`
	OrgName  string `yaml:"org_name" validate:"required"`
}

type SetcloneGitlab struct {
	RootPath  string `yaml:"root_path" validate:"required"`
	Provider  string `yaml:"provider" validate:"required"`
	Url       string `yaml:"url"`
	Recursive bool   `yaml:"recursive"`
	Protocol  string `yaml:"protocol" validate:"required,oneof=http https ssh"`
	GroupName string `yaml:"group_name" binding:"required"`
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

	err = cfg.validateConfig()
	if err != nil {
		printValidationErrors(err)
		log.Fatalf("failed to validate config file: %v", err)
	}
}

// 커스텀 오류 메시지 맵
var errorMessages = map[string]string{
	"required": "필수 필드입니다.",
	"url":      "유효한 URL을 입력해야 합니다.",
	"oneof":    "허용되는 값이 아닙니다 (http, https, ssh).",
}

// 유효성 검사 오류 상세 메시지 출력 함수
func printValidationErrors(err error) {
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		for _, e := range errs {
			// 기본 메시지
			msg, exists := errorMessages[e.Tag()]
			if !exists {
				msg = fmt.Sprintf("필드 '%s'은(는) '%s' 규칙을 만족해야 합니다.", e.Field(), e.Tag())
			}

			// 추가 정보가 필요한 경우 (예: oneof)
			if e.Tag() == "oneof" {
				msg = fmt.Sprintf("필드 '%s'은(는) 허용되는 값 중 하나여야 합니다: %s.", e.Field(), e.Param())
			}

			fmt.Printf("오류: %s\n", msg)
		}
	}
}

// 유효성 검사 함수
func (cfg *setcloneConfig) validateConfig() error {
	// validate := validator.New(validator.WithRequiredStructEnabled())
	validate := validator.New()
	return validate.Struct(cfg)
}
