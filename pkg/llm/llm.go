package llm

import (
	"context"
	"fmt"

	"github.com/soralabs/zen/pkg/logger"
)

// LLMClient is the main client that manages provider interactions
type LLMClient struct {
	provider Provider
	logger   *logger.Logger
	ctx      context.Context
}

// NewLLMClient creates a new LLM client with the specified provider
func NewLLMClient(config Config) (*LLMClient, error) {
	var provider Provider

	switch config.ProviderType {
	case ProviderOpenAI:
		provider = NewOpenAIProvider(config)
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", config.ProviderType)
	}

	return &LLMClient{
		provider: provider,
		logger:   config.Logger,
		ctx:      config.Context,
	}, nil
}

// Implement the Provider interface by delegating to the actual provider
func (c *LLMClient) GenerateCompletion(req CompletionRequest) (Message, error) {
	return c.provider.GenerateCompletion(c.ctx, req)
}

func (c *LLMClient) GenerateStructuredOutput(req StructuredOutputRequest, result interface{}) error {
	return c.provider.GenerateStructuredOutput(c.ctx, req, result)
}

func (c *LLMClient) EmbedText(text string) ([]float32, error) {
	return c.provider.EmbedText(c.ctx, text)
}

// Helper functions for creating messages
func NewSystemMessage(content string) Message {
	return Message{
		Role:    RoleSystem,
		Content: content,
	}
}

func NewUserMessage(content string) Message {
	return Message{
		Role:    RoleUser,
		Content: content,
	}
}

func NewAssistantMessage(content string) Message {
	return Message{
		Role:    RoleAssistant,
		Content: content,
	}
}

func NewToolMessage(content string, name string) Message {
	return Message{
		Role:    RoleTool,
		Content: content,
		Name:    name,
	}
}
