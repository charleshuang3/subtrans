package translator

import (
	"fmt"

	"github.com/charleshuang3/subtrans/pkg/config"
	"github.com/charleshuang3/subtrans/pkg/sub"
	"github.com/google/jsonschema-go/jsonschema"
)

const (
	promptTmpl = `Translate the following subtitle texts to %s. Return a JSON object with a "translations" array containing the translated texts in the same order:

[%s]

Return format:
{
  "translations": ["translation1", "translation2", ...]
}`
)

func NewLLMTranslator(cfg *config.Config) (sub.Translator, error) {
	if cfg.API == config.OpenAI {
		return newOpenAITranslator(cfg), nil
	} else if cfg.API == config.Gemini {
		return newGeminiTranslator(cfg)
	}

	return nil, fmt.Errorf("Unsupported API type: %s", cfg.API)
}

type TranslationResponse struct {
	Translations []string `json:"translations"`
}

var (
	translationResponseJSONSchema, _ = jsonschema.For[TranslationResponse](&jsonschema.ForOptions{})
)
