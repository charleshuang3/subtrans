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

func TestOpenAI(t *testing.T) {
	tests := []struct {
		name            string
		keyName         string
		apiURL          string
		model           string
		structureOutput string
	}{
		{
			name:            "JSONOutput",
			keyName:         "DEEPSEEK_KEY",
			apiURL:          "https://api.deepseek.com",
			model:           "deepseek-chat",
			structureOutput: config.OpenAIJSONObject,
		},
		{
			name:            "JSONSchema",
			keyName:         "XAI_KEY",
			apiURL:          "https://api.x.ai/v1",
			model:           "grok-4-1-fast-non-reasoning",
			structureOutput: config.OpenAIJSONSchema,
		},
	}

	testInput := []string{"hello", "The sky is blue"}
	wantOutput := []string{"你好", "天空是蓝色的"}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := getKey(t, tt.keyName)
			cfg := config.Config{
				API:             config.OpenAI,
				APIKey:          key,
				APIURL:          tt.apiURL,
				Model:           tt.model,
				StructureOutput: tt.structureOutput,
				TargetLang:      "简体中文",
			}
			translator := NewLLMTranslator(&cfg)
			got, err := translator.Translate(testInput)
			require.NoError(t, err)
			require.Equal(t, wantOutput, got)
		})
	}
}
