package translator

import (
	"github.com/charleshuang3/subtrans/pkg/config"
	"github.com/charleshuang3/subtrans/pkg/sub"
)

const (
	promptTmpl = `Translate the following subtitle texts to %s. Return a JSON object with a "translations" array containing the translated texts in the same order:

[%s]

Return format:
{
  "translations": ["translation1", "translation2", ...]
}`
)

func NewLLMTranslator(cfg *config.Config) sub.Translator {
	if cfg.API == config.OpenAI {
		return newOpenAITranslator(cfg)
	}
	return nil
}
