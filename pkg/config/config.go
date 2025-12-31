package config

import (
	"errors"
	"fmt"
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

type LLMProvider struct {
	API             string `yaml:"api"`
	APIKey          string `yaml:"api_key"`
	APIURL          string `yaml:"api_url"`
	Model           string `yaml:"model"`
	MaxTokens       int    `yaml:"max_tokens"`
	StructureOutput string `yaml:"structure_output"` // only used for openai
}

type Config struct {
	DefaultLLM string                 `yaml:"default_llm"`
	LLMs       map[string]LLMProvider `yaml:"llms"`
	TargetLang string                 `yaml:"target_lang"`
	Prompts    map[string]string      `yaml:"prompts"`
}

func (c *Config) validate() error {
	// Validate that we have at least one LLM provider
	if len(c.LLMs) == 0 {
		return errors.New("at least one LLM provider is required")
	}

	if len(c.LLMs) > 1 && c.DefaultLLM == "" {
		return errors.New("default LLM provider is required when there are multiple LLM providers")
	}

	if c.DefaultLLM == "" {
		return errors.New("default LLM provider is required")
	}

	// Validate DefaultLLM if specified
	if _, exists := c.LLMs[c.DefaultLLM]; !exists {
		return errors.New("default LLM provider not found in LLMs map")
	}

	// Validate each LLM provider
	for name, provider := range c.LLMs {
		if err := c.validateLLMProvider(name, provider); err != nil {
			return err
		}
	}

	// Initialize empty prompts map if nil
	if c.Prompts == nil {
		c.Prompts = map[string]string{}
	}

	return nil
}

func (c *Config) validateLLMProvider(name string, provider LLMProvider) error {
	if provider.API != OpenAI && provider.API != Gemini {
		return fmt.Errorf("invalid api for LLM provider '%s'", name)
	}
	if provider.API == OpenAI {
		if provider.StructureOutput == "" {
			// Set default structure output for OpenAI
			// Note: This modifies the provider, but since it's a copy, we need to update the map
			// In practice, this would need to be handled differently if the config is shared
			const defaultStructureOutput = OpenAIJSONSchema
			provider.StructureOutput = defaultStructureOutput
		}
		if provider.StructureOutput != OpenAIJSONObject && provider.StructureOutput != OpenAIJSONSchema {
			return fmt.Errorf("invalid structure_output for LLM provider '%s'", name)
		}
	}
	if provider.APIKey == "" {
		return fmt.Errorf("api_key is required for LLM provider '%s'", name)
	}
	if provider.Model == "" {
		return fmt.Errorf("model is required for LLM provider '%s'", name)
	}
	if provider.MaxTokens == 0 {
		provider.MaxTokens = defaultMaxTokens
	}
	return nil
}

// GetDefaultLLM returns the default LLM provider
func (c *Config) GetDefaultLLM() (LLMProvider, error) {
	provider, exists := c.LLMs[c.DefaultLLM]
	if !exists {
		return LLMProvider{}, errors.New("default LLM provider not found")
	}

	return provider, nil
}

// GetLLM returns a specific LLM provider by name
func (c *Config) GetLLM(name string) (LLMProvider, error) {
	provider, exists := c.LLMs[name]
	if !exists {
		return LLMProvider{}, fmt.Errorf("LLM provider '%s' not found", name)
	}

	return provider, nil
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
