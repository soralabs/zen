package twitter

import (
	"context"
	"fmt"
	"time"

	"github.com/soralabs/zen/llm"
	"github.com/soralabs/zen/logger"
	"github.com/soralabs/zen/options"

	"gorm.io/gorm"
)

// ValidateRequiredFields checks if all required dependencies are properly initialized.
// Returns an error if any required field is nil.
func (k *Twitter) ValidateRequiredFields() error {
	if k.ctx == nil {
		return fmt.Errorf("context is required")
	}
	if k.logger == nil {
		return fmt.Errorf("logger is required")
	}
	if k.database == nil {
		return fmt.Errorf("database is required")
	}
	if k.llmClient == nil {
		return fmt.Errorf("LLM client is required")
	}
	return nil
}

// WithContext sets the context for the Twitter client.
// The context is used for cancellation and timeout control.
func WithContext(ctx context.Context) options.Option[Twitter] {
	return func(k *Twitter) error {
		k.ctx = ctx
		return nil
	}
}

// WithLogger sets the logger instance for the Twitter client.
// The logger is used for recording operational events and errors.
func WithLogger(logger *logger.Logger) options.Option[Twitter] {
	return func(k *Twitter) error {
		k.logger = logger
		return nil
	}
}

// WithDatabase sets the database connection for the Twitter client.
// The database is used for persisting Twitter-related data.
func WithDatabase(database *gorm.DB) options.Option[Twitter] {
	return func(k *Twitter) error {
		k.database = database
		return nil
	}
}

// WithLLM sets the Language Learning Model client for the Twitter client.
// The LLM client is used for processing and generating Twitter content.
func WithLLM(llmClient *llm.LLMClient) options.Option[Twitter] {
	return func(k *Twitter) error {
		k.llmClient = llmClient
		return nil
	}
}

// WithTwitterMonitorInterval sets the minimum and maximum monitoring interval durations.
// Returns an error if min is greater than max.
// These intervals control how frequently the Twitter client checks for updates.
func WithTwitterMonitorInterval(min, max time.Duration) options.Option[Twitter] {
	return func(k *Twitter) error {
		if min > max {
			return fmt.Errorf("minimum interval cannot be greater than maximum interval")
		}
		k.twitterConfig.MonitorInterval = IntervalConfig{
			Min: min,
			Max: max,
		}
		return nil
	}
}

// WithTwitterCredentials sets the authentication credentials for the Twitter client.
// ct0: Twitter's ct0 cookie value
// authToken: Twitter's authentication token
// user: Twitter username
func WithTwitterCredentials(ct0, authToken, user string) options.Option[Twitter] {
	return func(k *Twitter) error {
		k.twitterConfig.Credentials.CT0 = ct0
		k.twitterConfig.Credentials.AuthToken = authToken
		k.twitterConfig.Credentials.User = user
		return nil
	}
}
