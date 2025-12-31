package translator

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charleshuang3/subtrans/pkg/config"
	"github.com/charleshuang3/subtrans/pkg/sub"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/tiktoken-go/tokenizer"
)

const (
	defaultPromptTmpl = `Translate the following subtitle texts to $TARGET_LANG$. Return a JSON object with a "translations" array containing the translated texts in the same order:

Return format:
{
  "translations": ["translation1", "translation2", ...]
}
  
Subtitle texts:
$SUBTITLES$
`
)

func NewLLMTranslator(cfg *config.Config, promptKey, llmProvider string) (sub.Translator, error) {
	var provider config.LLMProvider
	var err error

	if llmProvider == "default" {
		provider, err = cfg.GetDefaultLLM()
		if err != nil {
			return nil, fmt.Errorf("failed to get default LLM provider: %w", err)
		}
	} else {
		provider, err = cfg.GetLLM(llmProvider)
		if err != nil {
			return nil, fmt.Errorf("failed to get LLM provider '%s': %w", llmProvider, err)
		}
	}

	switch provider.API {
	case config.OpenAI:
		return newOpenAITranslator(cfg, provider, promptKey), nil
	case config.Gemini:
		return newGeminiTranslator(cfg, provider, promptKey)
	default:
		return nil, fmt.Errorf("Unsupported API type: %s", provider.API)
	}
}

func getPromptTmpl(cfg *config.Config, promptKey string) (string, error) {
	if s, ok := cfg.Prompts[promptKey]; ok {
		return s, nil
	}
	if promptKey == "default" {
		return defaultPromptTmpl, nil
	}
	return "", fmt.Errorf("prompt %q not found from config", promptKey)
}

func toPrompt(promptTmpl string, lang string, texts []string) (string, error) {
	textsJSON, err := json.Marshal(texts)
	if err != nil {
		return "", fmt.Errorf("failed to marshal input texts: %w", err)
	}

	s := strings.ReplaceAll(promptTmpl, "$TARGET_LANG$", lang)
	return strings.ReplaceAll(s, "$SUBTITLES$", string(textsJSON)), nil
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
