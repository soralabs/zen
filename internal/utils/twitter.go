package utils

import (
	"fmt"
	"time"

	"github.com/soralabs/zen/pkg/twitter"

	"github.com/soralabs/zen/db"
	"github.com/soralabs/zen/id"

	"github.com/mitchellh/mapstructure"
	"github.com/pgvector/pgvector-go"
)

// Helper method to create fragment from tweet
func CreateTweetFragment(tweet *twitter.ParsedTweet, actorId id.ID, embedding []float32) (*db.Fragment, error) {
	var metadata db.Metadata
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  &metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create decoder: %w", err)
	}
	if err := decoder.Decode(tweet); err != nil {
		return nil, fmt.Errorf("failed to decode tweet metadata: %w", err)
	}

	return &db.Fragment{
		ID:        id.FromString(tweet.TweetID),
		ActorID:   actorId,
		SessionID: id.FromString(tweet.TweetConversationID),
		Content:   tweet.TweetText,
		Embedding: pgvector.NewVector(embedding),
		Metadata:  metadata,
		CreatedAt: time.Unix(tweet.TweetCreatedAt, 0),
	}, nil
}

// DecodeTweetMetadata is a helper function to decode tweet metadata using json tags
func DecodeTweetMetadata(input interface{}, result *twitter.ParsedTweet) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  result,
	})
	if err != nil {
		return fmt.Errorf("failed to create decoder: %w", err)
	}

	if err := decoder.Decode(input); err != nil {
		return fmt.Errorf("failed to decode tweet metadata: %w", err)
	}

	return nil
}
