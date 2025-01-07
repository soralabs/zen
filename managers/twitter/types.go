package twitter_manager

import (
	"time"

	"github.com/soralabs/zen/manager"
	"github.com/soralabs/zen/options"
	"github.com/soralabs/zen/pkg/twitter"
	"github.com/soralabs/zen/state"
)

// Package twitter provides functionality for managing Twitter-specific interactions
// and conversation handling for the agent system.

// StateDataKey constants for different types of Twitter data
const (
	// TwitterConversations represents formatted conversation threads
	TwitterConversations state.StateDataKey = "twitter_conversations"
)

// TwitterManager handles Twitter-specific functionality and conversation management
type TwitterManager struct {
	*manager.BaseManager
	options.RequiredFields

	twitterClient   *twitter.Client // Client for Twitter API interactions
	twitterUsername string
}

// TwitterActions represents the possible actions an agent can take on Twitter
type TwitterActions struct {
	ShouldTweet bool   // Whether to send a tweet response
	ShouldLike  bool   // Whether to like the tweet
	Reasoning   string // Explanation for the chosen actions
}

// ParsedTweet represents a processed Twitter message with extracted metadata
type ParsedTweet struct {
	TweetID             string    // Unique identifier for the tweet
	UserID              string    // ID of the tweet author
	UserName            string    // Username of the tweet author
	TweetText           string    // Content of the tweet
	InReplyToTweetID    string    // ID of the tweet this is replying to
	TweetConversationID string    // ID of the root conversation
	CreatedAt           time.Time // When the tweet was created
}
