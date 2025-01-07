package llm

import (
	"context"

	"github.com/soralabs/zen/logger"
)

// ProviderType identifies different LLM providers
type ProviderType string

const (
	ProviderOpenAI ProviderType = "openai"
)

// Config holds the configuration for an LLM provider
type Config struct {
	ProviderType ProviderType
	APIKey       string
	ModelConfig  map[ModelType]string // Maps capability levels to specific model names
	Logger       *logger.Logger
	Context      context.Context
}
