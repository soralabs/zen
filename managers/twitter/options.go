package twitter_manager

import (
	"fmt"

	"github.com/soralabs/zen/options"
	"github.com/soralabs/zen/pkg/twitter"
)

func (t *TwitterManager) ValidateRequiredFields() error {
	if t.twitterClient == nil {
		return fmt.Errorf("twitter client is required")
	}
	if t.twitterUsername == "" {
		return fmt.Errorf("twitter username is required")
	}
	return nil
}

func WithTwitterClient(client *twitter.Client) options.Option[TwitterManager] {
	return func(m *TwitterManager) error {
		m.twitterClient = client
		return nil
	}
}

func WithTwitterUsername(username string) options.Option[TwitterManager] {
	return func(m *TwitterManager) error {
		m.twitterUsername = username
		return nil
	}
}
