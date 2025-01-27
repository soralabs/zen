package llm

import (
	"context"

	"github.com/soralabs/zen/logger"
)

// ProviderType identifies different LLM providers
type ProviderType string

const (
	ProviderOpenAI   ProviderType = "openai"
	ProviderDeepseek ProviderType = "deepseek"
)

// ProviderConfig holds the configuration for a specific provider
type ProviderConfig struct {
	Type        ProviderType
	APIKey      string
	ModelConfig map[ModelType]string // Maps capability levels to specific model names
}

// Config holds the configuration for the LLM client
type Config struct {
	DefaultProvider ProviderConfig
	// Specific providers for different capabilities
	EmbeddingProvider *ProviderConfig // If nil, uses DefaultProvider
	ChatProvider      *ProviderConfig // If nil, uses DefaultProvider
	Logger            *logger.Logger
	Context           context.Context
}
