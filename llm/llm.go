package llm

import (
	"context"
	"fmt"

	"github.com/soralabs/zen/logger"
)

// LLMClient is the main client that manages provider interactions
type LLMClient struct {
	defaultProvider   Provider
	chatProvider      Provider
	embeddingProvider Provider
	logger            *logger.Logger
	ctx               context.Context
}

func createProvider(config ProviderConfig, logger *logger.Logger, ctx context.Context) (Provider, error) {
	switch config.Type {
	case ProviderOpenAI:
		return NewOpenAIProvider(Config{
			DefaultProvider: config,
			Logger:          logger,
			Context:         ctx,
		}), nil
	case ProviderDeepseek:
		return NewDeepseekProvider(Config{
			DefaultProvider: config,
			Logger:          logger,
			Context:         ctx,
		}), nil
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", config.Type)
	}
}

// NewLLMClient creates a new LLM client with the specified providers
func NewLLMClient(config Config) (*LLMClient, error) {
	defaultProvider, err := createProvider(config.DefaultProvider, config.Logger, config.Context)
	if err != nil {
		return nil, fmt.Errorf("failed to create default provider: %w", err)
	}

	var chatProvider Provider = defaultProvider
	if config.ChatProvider != nil {
		chatProvider, err = createProvider(*config.ChatProvider, config.Logger, config.Context)
		if err != nil {
			return nil, fmt.Errorf("failed to create chat provider: %w", err)
		}
	}

	var embeddingProvider Provider = defaultProvider
	if config.EmbeddingProvider != nil {
		embeddingProvider, err = createProvider(*config.EmbeddingProvider, config.Logger, config.Context)
		if err != nil {
			return nil, fmt.Errorf("failed to create embedding provider: %w", err)
		}
	}

	return &LLMClient{
		defaultProvider:   defaultProvider,
		chatProvider:      chatProvider,
		embeddingProvider: embeddingProvider,
		logger:            config.Logger,
		ctx:               config.Context,
	}, nil
}

// Implement the Provider interface by delegating to the appropriate provider
func (c *LLMClient) GenerateCompletion(req CompletionRequest) (Message, error) {
	return c.chatProvider.GenerateCompletion(c.ctx, req)
}

func (c *LLMClient) GenerateStructuredOutput(req StructuredOutputRequest, result interface{}) error {
	return c.chatProvider.GenerateStructuredOutput(c.ctx, req, result)
}

func (c *LLMClient) EmbedText(text string) ([]float32, error) {
	return c.embeddingProvider.EmbedText(c.ctx, text)
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
