package repowatch

import (
	"os"

	"gopkg.in/yaml.v2"
)

const DEFAULT_CONFI_FILE_PATH = "./config.yaml"

type Github struct {
	Token string `yaml:"token"`
}

type LlmConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Features struct {
	Summary bool `yaml:"summary"`
}

type RepoDef struct {
	Name      string `yaml:"name"`
	LimitNum  int    `yaml:"limitNum"`
	LimitDays int    `yaml:"limitDays"`
}


type Config struct {
	Github      Github `yaml:"github"`
	Features    Features `yaml:"features"`
	LlmConfig   LlmConfig `yaml:"llm"`
	Repos       []RepoDef `yaml:"repos"`
}

func ParseConfig() (*Config, error) {
	filepath := os.Getenv("CONFIG_FILE_PATH")
	if filepath == "" {
		filepath = DEFAULT_CONFI_FILE_PATH
	}

	configContents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(configContents, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
