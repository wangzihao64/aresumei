package llm

import (
	"context"
	"errors"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

type OpenAIConfig struct {
	APIKey  string
	Model   string
	BaseURL string
	Timeout time.Duration
}

func NewOpenAIChatModel(ctx context.Context, cfg OpenAIConfig) (model.BaseChatModel, error) {
	if cfg.APIKey == "" {
		return nil, errors.New("openai api key is required")
	}
	if cfg.Model == "" {
		return nil, errors.New("openai model is required")
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 60 * time.Second
	}

	return openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  cfg.APIKey,
		Model:   cfg.Model,
		BaseURL: cfg.BaseURL,
		Timeout: cfg.Timeout,
	})
}
