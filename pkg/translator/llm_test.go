package translator

import (
	"os"
	"testing"

	"github.com/charleshuang3/subtrans/pkg/config"
	"github.com/stretchr/testify/require"
)

func getKey(t *testing.T, name string) string {
	t.Helper()

	key := os.Getenv(name)
	if key == "" {
		t.Skipf("%s not set", name)
	}
	return key
}

func TestLLMTranslator(t *testing.T) {
	tests := []struct {
		api             string
		name            string
		keyName         string
		apiURL          string
		model           string
		structureOutput string
	}{
		{
			api:             "openai",
			name:            "OpenAI_JSONOutput",
			keyName:         "DEEPSEEK_KEY",
			apiURL:          "https://api.deepseek.com",
			model:           "deepseek-chat",
			structureOutput: config.OpenAIJSONObject,
		},
		{
			api:             "openai",
			name:            "OpenAI_JSONSchema",
			keyName:         "XAI_KEY",
			apiURL:          "https://api.x.ai/v1",
			model:           "grok-4-1-fast-non-reasoning",
			structureOutput: config.OpenAIJSONSchema,
		},
		{
			api:     "gemini",
			name:    "Gemini",
			keyName: "GEMINI_KEY",
			model:   "gemini-2.5-flash-lite",
		},
	}

	testInput := []string{"hello", "The sky is blue"}
	wantOutput := []string{"你好", "天空是蓝色的"}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := getKey(t, tt.keyName)
			cfg := config.Config{
				DefaultLLM: "test-provider",
				LLMs: map[string]config.LLMProvider{
					"test-provider": {
						API:             tt.api,
						APIKey:          key,
						APIURL:          tt.apiURL,
						Model:           tt.model,
						StructureOutput: tt.structureOutput,
					},
				},
				TargetLang: "简体中文",
			}
			translator, err := NewLLMTranslator(&cfg, "default", "default")
			require.NoError(t, err)
			got, err := translator.Translate(testInput)
			require.NoError(t, err)
			require.Equal(t, wantOutput, got)
		})
	}
}

func TestGetPromptTmpl(t *testing.T) {
	tests := []struct {
		name      string
		config    config.Config
		promptKey string
		want      string
		wantError bool
	}{
		{
			name: "Custom prompt from config",
			config: config.Config{
				Prompts: map[string]string{
					"custom": "Custom prompt template for $TARGET_LANG$",
				},
			},
			promptKey: "custom",
			want:      "Custom prompt template for $TARGET_LANG$",
			wantError: false,
		},
		{
			name:      "Default prompt",
			config:    config.Config{Prompts: map[string]string{}},
			promptKey: "default",
			want:      defaultPromptTmpl,
			wantError: false,
		},
		{
			name: "Non-existent prompt key",
			config: config.Config{
				Prompts: map[string]string{
					"existing": "Some prompt",
				},
			},
			promptKey: "nonexistent",
			want:      "",
			wantError: true,
		},
		{
			name:      "Empty prompts map with default key",
			config:    config.Config{Prompts: map[string]string{}},
			promptKey: "default",
			want:      defaultPromptTmpl,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPromptTmpl(&tt.config, tt.promptKey)
			if tt.wantError {
				require.Error(t, err)
				require.Empty(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestToPrompt(t *testing.T) {
	tests := []struct {
		name       string
		promptTmpl string
		lang       string
		texts      []string
		want       string
	}{
		{
			name:       "Basic template substitution",
			promptTmpl: "Translate to $TARGET_LANG$: $SUBTITLES$",
			lang:       "Spanish",
			texts:      []string{"Hello", "How are you?"},
			want:       "Translate to Spanish: [\"Hello\",\"How are you?\"]",
		},
		{
			name:       "Default prompt template",
			promptTmpl: defaultPromptTmpl,
			lang:       "French",
			texts:      []string{"Good morning", "See you later"},
			want: `Translate the following subtitle texts to French. Return a JSON object with a "translations" array containing the translated texts in the same order:

Return format:
{
  "translations": ["translation1", "translation2", ...]
}
  
Subtitle texts:
["Good morning","See you later"]
`,
		},
		{
			name:       "Empty texts array",
			promptTmpl: "Translate to $TARGET_LANG$: $SUBTITLES$",
			lang:       "German",
			texts:      []string{},
			want:       "Translate to German: []",
		},
		{
			name:       "Single text",
			promptTmpl: "$TARGET_LANG$ translation: $SUBTITLES$",
			lang:       "Japanese",
			texts:      []string{"こんにちは"},
			want:       "Japanese translation: [\"こんにちは\"]",
		},
		{
			name:       "Text with special characters",
			promptTmpl: "Translate to $TARGET_LANG$: $SUBTITLES$",
			lang:       "Chinese",
			texts:      []string{"Hello\nworld", "Tab\there", "Quote\"test"},
			want:       "Translate to Chinese: [\"Hello\\nworld\",\"Tab\\there\",\"Quote\\\"test\"]",
		},
		{
			name:       "Template without target lang placeholder",
			promptTmpl: "Translate: $SUBTITLES$",
			lang:       "English",
			texts:      []string{"Hello"},
			want:       "Translate: [\"Hello\"]",
		},
		{
			name:       "Template without subtitles placeholder",
			promptTmpl: "Translate to $TARGET_LANG$",
			lang:       "Italian",
			texts:      []string{"Hello"},
			want:       "Translate to Italian",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toPrompt(tt.promptTmpl, tt.lang, tt.texts)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
