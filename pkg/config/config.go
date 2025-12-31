package config

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

const (
	OpenAI           = "openai"
	Gemini           = "gemini"
	OpenAIJSONObject = "json_object"
	OpenAIJSONSchema = "json_schema"
	defaultMaxTokens = 128000 // llm usually works better on small context
)

type Config struct {
	API             string            `yaml:"api"`
	APIKey          string            `yaml:"api_key"`
	APIURL          string            `yaml:"api_url"`
	Model           string            `yaml:"model"`
	MaxTokens       int               `yaml:"max_tokens"`
	StructureOutput string            `yaml:"structure_output"` // only used for openai
	TargetLang      string            `yaml:"target_lang"`
	Prompts         map[string]string `yaml:"prompts"`
}

func (c *Config) validate() error {
	if c.MaxTokens == 0 {
		c.MaxTokens = defaultMaxTokens
	}
	if c.API != OpenAI && c.API != Gemini {
		return errors.New("invalid api")
	}
	if c.API == OpenAI {
		if c.StructureOutput == "" {
			c.StructureOutput = OpenAIJSONSchema
		}
		if c.StructureOutput != OpenAIJSONObject && c.StructureOutput != OpenAIJSONSchema {
			return errors.New("invalid structure_output")
		}
	}
	if c.APIKey == "" {
		return errors.New("api_key is required")
	}
	if c.Model == "" {
		return errors.New("model is required")
	}
	if c.Prompts == nil {
		c.Prompts = map[string]string{}
	}
	return nil
}

func Read(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func FindConfig(path string) (string, error) {
	paths := []string{}

	if path != "" {
		paths = append(paths, path)
	}
	paths = append(paths, ".env.yaml")

	usr, err := user.Current()
	if err == nil {
		paths = append(paths, filepath.Join(usr.HomeDir, ".config", "subtrans", "config.yaml"))
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", errors.New("config file not found")
}
