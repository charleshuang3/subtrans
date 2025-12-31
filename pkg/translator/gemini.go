package translator

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/charleshuang3/subtrans/pkg/config"
	"google.golang.org/genai"
)

type GeminiTranslator struct {
	Config *config.Config
	client *genai.Client
}

func newGeminiTranslator(cfg *config.Config) (*GeminiTranslator, error) {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  cfg.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Gemini client: %w", err)
	}

	return &GeminiTranslator{
		Config: cfg,
		client: client,
	}, nil
}

func (t *GeminiTranslator) Length(text string) int {
	return tokenCount(text)
}

func (t *GeminiTranslator) MaxLength() int {
	return int(float64(t.Config.MaxTokens) * 0.95)
}

func (t *GeminiTranslator) Translate(texts []string) ([]string, error) {
	if len(texts) == 0 {
		return []string{}, nil
	}

	textsJSON, err := json.Marshal(texts)
	if err != nil {
		return texts, fmt.Errorf("failed to marshal input texts: %w", err)
	}

	prompt := fmt.Sprintf(promptTmpl, t.Config.TargetLang, string(textsJSON))

	ctx := context.Background()

	generateConfig := &genai.GenerateContentConfig{
		ResponseMIMEType:   "application/json",
		ResponseJsonSchema: translationResponseJSONSchema,
	}

	resp, err := t.client.Models.GenerateContent(ctx, t.Config.Model, genai.Text(prompt), generateConfig)
	if err != nil {
		return texts, err
	}

	if len(resp.Candidates) == 0 {
		return texts, fmt.Errorf("no completion choices returned from Gemini API")
	}

	content := resp.Candidates[0].Content.Parts[0]
	if content.Text == "" {
		return texts, fmt.Errorf("empty response from Gemini API")
	}

	var result TranslationResponse
	if err := json.Unmarshal([]byte(content.Text), &result); err != nil {
		return texts, fmt.Errorf("failed to unmarshal translation response: %w", err)
	}

	if len(result.Translations) == len(texts) {
		return result.Translations, nil
	}

	return texts, fmt.Errorf("translation count mismatch: got %d translations for %d input texts", len(result.Translations), len(texts))
}
