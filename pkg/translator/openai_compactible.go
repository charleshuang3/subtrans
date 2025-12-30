package translator

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/charleshuang3/subtrans/pkg/config"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

type OpenAICompactibleTranslator struct {
	Config *config.Config
	client openai.Client
}

type TranslationResponse struct {
	Translations []string `json:"translations"`
}

var (
	translationResponseJSONSchema, _ = jsonschema.For[TranslationResponse](&jsonschema.ForOptions{})
)

func newOpenAITranslator(cfg *config.Config) *OpenAICompactibleTranslator {
	client := openai.NewClient(
		option.WithAPIKey(cfg.APIKey),
		option.WithBaseURL(cfg.APIURL),
	)

	return &OpenAICompactibleTranslator{
		Config: cfg,
		client: client,
	}
}

func (t *OpenAICompactibleTranslator) Translate(texts []string) ([]string, error) {
	if len(texts) == 0 {
		return []string{}, nil
	}

	textsJSON, err := json.Marshal(texts)
	if err != nil {
		return texts, fmt.Errorf("failed to marshal input texts: %w", err)
	}

	prompt := fmt.Sprintf(promptTmpl, t.Config.TargetLang, string(textsJSON))

	ctx := context.Background()

	responseFormat := openai.ChatCompletionNewParamsResponseFormatUnion{}
	if t.Config.StructureOutput == config.OpenAIJSONObject {
		param := shared.NewResponseFormatJSONObjectParam()
		responseFormat.OfJSONObject = &param
	} else {
		responseFormat.OfJSONSchema = &shared.ResponseFormatJSONSchemaParam{
			JSONSchema: shared.ResponseFormatJSONSchemaJSONSchemaParam{
				Name:   "translation_response",
				Schema: translationResponseJSONSchema,
			},
			Type: "json_schema",
		}
	}

	completion, err := t.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: shared.ChatModel(t.Config.Model),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		ResponseFormat: responseFormat,
	})

	if err != nil {
		return texts, fmt.Errorf("failed to get completion from OpenAI API: %w", err)
	}

	if len(completion.Choices) == 0 {
		return texts, fmt.Errorf("no completion choices returned from OpenAI API")
	}

	content := completion.Choices[0].Message.Content

	var result TranslationResponse
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return texts, fmt.Errorf("failed to unmarshal translation response: %w", err)
	}

	if len(result.Translations) == len(texts) {
		return result.Translations, nil
	}

	return texts, fmt.Errorf("translation count mismatch: got %d translations for %d input texts", len(result.Translations), len(texts))
}
