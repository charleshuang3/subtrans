package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr string
	}{
		{
			name:    "no LLM providers",
			config:  Config{},
			wantErr: "at least one LLM provider is required",
		},
		{
			name: "multiple LLMs without default",
			config: Config{
				LLMs: map[string]LLMProvider{
					"openai": {API: OpenAI, APIKey: "key", Model: "gpt-4"},
					"gemini": {API: Gemini, APIKey: "key", Model: "gemini-pro"},
				},
			},
			wantErr: "default LLM provider is required when there are multiple LLM providers",
		},
		{
			name: "single LLM without default",
			config: Config{
				LLMs: map[string]LLMProvider{
					"openai": {API: OpenAI, APIKey: "key", Model: "gpt-4"},
				},
			},
			wantErr: "default LLM provider is required",
		},
		{
			name: "default LLM not found",
			config: Config{
				DefaultLLM: "nonexistent",
				LLMs: map[string]LLMProvider{
					"openai": {API: OpenAI, APIKey: "key", Model: "gpt-4"},
				},
			},
			wantErr: "default LLM provider not found in LLMs map",
		},
		{
			name: "valid config",
			config: Config{
				DefaultLLM: "openai",
				LLMs: map[string]LLMProvider{
					"openai": {API: OpenAI, APIKey: "key", Model: "gpt-4"},
				},
			},
			wantErr: "",
		},
		{
			name: "valid config initializes nil prompts",
			config: Config{
				DefaultLLM: "openai",
				LLMs: map[string]LLMProvider{
					"openai": {API: OpenAI, APIKey: "key", Model: "gpt-4"},
				},
				Prompts: nil,
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()
			if tt.wantErr == "" {
				assert.NoError(t, err)
				assert.NotNil(t, tt.config.Prompts, "validate() should initialize nil Prompts map")
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}

func TestConfig_validateLLMProvider(t *testing.T) {
	tests := []struct {
		name     string
		llmName  string
		provider LLMProvider
		wantErr  string
	}{
		{
			name:     "invalid API type",
			llmName:  "test",
			provider: LLMProvider{API: "invalid", APIKey: "key", Model: "model"},
			wantErr:  "invalid api for LLM provider 'test'",
		},
		{
			name:     "missing API key",
			llmName:  "test",
			provider: LLMProvider{API: OpenAI, Model: "gpt-4"},
			wantErr:  "api_key is required for LLM provider 'test'",
		},
		{
			name:     "missing model",
			llmName:  "test",
			provider: LLMProvider{API: OpenAI, APIKey: "key"},
			wantErr:  "model is required for LLM provider 'test'",
		},
		{
			name:     "invalid structure_output for OpenAI",
			llmName:  "test",
			provider: LLMProvider{API: OpenAI, APIKey: "key", Model: "gpt-4", StructureOutput: "invalid"},
			wantErr:  "invalid structure_output for LLM provider 'test'",
		},
		{
			name:     "valid OpenAI provider with json_object",
			llmName:  "test",
			provider: LLMProvider{API: OpenAI, APIKey: "key", Model: "gpt-4", StructureOutput: OpenAIJSONObject},
			wantErr:  "",
		},
		{
			name:     "valid OpenAI provider with json_schema",
			llmName:  "test",
			provider: LLMProvider{API: OpenAI, APIKey: "key", Model: "gpt-4", StructureOutput: OpenAIJSONSchema},
			wantErr:  "",
		},
		{
			name:     "valid OpenAI provider with default structure_output",
			llmName:  "test",
			provider: LLMProvider{API: OpenAI, APIKey: "key", Model: "gpt-4"},
			wantErr:  "",
		},
		{
			name:     "valid Gemini provider",
			llmName:  "test",
			provider: LLMProvider{API: Gemini, APIKey: "key", Model: "gemini-pro"},
			wantErr:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{}
			err := c.validateLLMProvider(tt.llmName, tt.provider)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}

func TestConfig_GetDefaultLLM(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		want    LLMProvider
		wantErr bool
	}{
		{
			name: "default LLM exists",
			config: Config{
				DefaultLLM: "openai",
				LLMs: map[string]LLMProvider{
					"openai": {API: OpenAI, APIKey: "key", Model: "gpt-4"},
				},
			},
			want:    LLMProvider{API: OpenAI, APIKey: "key", Model: "gpt-4"},
			wantErr: false,
		},
		{
			name: "default LLM not found",
			config: Config{
				DefaultLLM: "nonexistent",
				LLMs:       map[string]LLMProvider{},
			},
			want:    LLMProvider{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.config.GetDefaultLLM()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestConfig_GetLLM(t *testing.T) {
	config := Config{
		LLMs: map[string]LLMProvider{
			"openai": {API: OpenAI, APIKey: "key", Model: "gpt-4"},
			"gemini": {API: Gemini, APIKey: "key2", Model: "gemini-pro"},
		},
	}

	tests := []struct {
		name    string
		llmName string
		want    LLMProvider
		wantErr bool
	}{
		{
			name:    "existing LLM",
			llmName: "openai",
			want:    LLMProvider{API: OpenAI, APIKey: "key", Model: "gpt-4"},
			wantErr: false,
		},
		{
			name:    "another existing LLM",
			llmName: "gemini",
			want:    LLMProvider{API: Gemini, APIKey: "key2", Model: "gemini-pro"},
			wantErr: false,
		},
		{
			name:    "nonexistent LLM",
			llmName: "nonexistent",
			want:    LLMProvider{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := config.GetLLM(tt.llmName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestRead(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "config_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name     string
		content  string
		wantErr  bool
		validate func(*testing.T, *Config)
	}{
		{
			name: "valid config",
			content: `default_llm: openai
llms:
  openai:
    api: openai
    api_key: test-key
    model: gpt-4
target_lang: en
`,
			wantErr: false,
			validate: func(t *testing.T, c *Config) {
				assert.Equal(t, "openai", c.DefaultLLM)
				assert.Equal(t, "en", c.TargetLang)
			},
		},
		{
			name:    "invalid yaml",
			content: `invalid: [yaml`,
			wantErr: true,
		},
		{
			name: "validation fails - no LLMs",
			content: `default_llm: openai
target_lang: en
`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write test config file
			configPath := filepath.Join(tmpDir, tt.name+".yaml")
			require.NoError(t, os.WriteFile(configPath, []byte(tt.content), 0644))

			cfg, err := Read(configPath)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, cfg)
				}
			}
		})
	}

	// Test file not found
	t.Run("file not found", func(t *testing.T) {
		_, err := Read("/nonexistent/path/config.yaml")
		assert.Error(t, err)
	})
}

func TestFindConfig(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "config_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a test config file
	testConfigPath := filepath.Join(tmpDir, "test-config.yaml")
	require.NoError(t, os.WriteFile(testConfigPath, []byte("test"), 0644))

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "explicit path exists",
			path:    testConfigPath,
			wantErr: false,
		},
		{
			name:    "explicit path does not exist",
			path:    "/nonexistent/config.yaml",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindConfig(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.path, got)
			}
		})
	}

	// Test empty path (should search default locations)
	t.Run("empty path searches defaults", func(t *testing.T) {
		// This will likely fail unless .env.yaml or ~/.config/subtrans/config.yaml exists
		_, err := FindConfig("")
		// We just verify it doesn't panic; the result depends on the environment
		_ = err
	})
}
