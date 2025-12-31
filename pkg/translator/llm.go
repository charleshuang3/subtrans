package translator

import (
	"fmt"

	"github.com/charleshuang3/subtrans/pkg/config"
	"github.com/charleshuang3/subtrans/pkg/sub"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/tiktoken-go/tokenizer"
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
	switch cfg.API {
	case config.OpenAI:
		return newOpenAITranslator(cfg), nil
	case config.Gemini:
		return newGeminiTranslator(cfg)
	default:
		return nil, fmt.Errorf("Unsupported API type: %s", cfg.API)
	}
}

type TranslationResponse struct {
	Translations []string `json:"translations"`
}

var (
	translationResponseJSONSchema, _ = jsonschema.For[TranslationResponse](&jsonschema.ForOptions{})
	encoder, _                       = tokenizer.Get(tokenizer.Cl100kBase)
)

// tokenCount is approximate: not all LLMs use cl100k, and tokenCount(a)+tokenCount(b) != tokenCount(a+b).
func tokenCount(sentence string) int {
	ids, _, _ := encoder.Encode(sentence)
	return len(ids)
}
