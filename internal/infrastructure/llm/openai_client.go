package llm

import (
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type Config struct {
	APIKey  string
	BaseURL string
}
type OpenAiClient struct {
	Client *openai.Client
}

func NewOpenAIClient(cfg Config) *OpenAiClient {
	client := openai.NewClient(option.WithAPIKey(cfg.APIKey),
		option.WithBaseURL(cfg.BaseURL))
	return &OpenAiClient{Client: &client}
}
