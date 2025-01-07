package twitter

import (
	"context"
	"time"

	"github.com/soralabs/zen/engine"
	"github.com/soralabs/zen/llm"
	"github.com/soralabs/zen/logger"
	"github.com/soralabs/zen/options"
	"github.com/soralabs/zen/pkg/twitter"

	"gorm.io/gorm"
)

type Twitter struct {
	options.RequiredFields

	ctx       context.Context
	logger    *logger.Logger
	database  *gorm.DB
	llmClient *llm.LLMClient

	assistant *engine.Engine

	twitterClient *twitter.Client
	twitterConfig TwitterConfig

	stopChan chan struct{}
}

type TwitterCredentials struct {
	CT0       string
	AuthToken string
	User      string
}

type IntervalConfig struct {
	Min time.Duration
	Max time.Duration
}

type TwitterConfig struct {
	MonitorInterval IntervalConfig
	Credentials     TwitterCredentials
}
