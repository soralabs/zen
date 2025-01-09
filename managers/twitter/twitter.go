package twitter_manager

import (
	"fmt"

	"github.com/soralabs/zen/db"
	"github.com/soralabs/zen/id"
	"github.com/soralabs/zen/internal/utils"
	"github.com/soralabs/zen/manager"
	"github.com/soralabs/zen/options"
	"github.com/soralabs/zen/pkg/twitter"
	"github.com/soralabs/zen/state"
)

// Package twitter_manager provides functionality for managing Twitter-specific interactions
// and conversation handling for the agent system.
// Does not process regular tweeting, but rather interactions like replies and likes.

// NewTwitterManager creates a new instance of TwitterManager with the provided options.
// Initializes the base manager and Twitter-specific configurations.
func NewTwitterManager(
	baseOpts []options.Option[manager.BaseManager],
	twitterOpts ...options.Option[TwitterManager],
) (*TwitterManager, error) {
	base, err := manager.NewBaseManager(baseOpts...)
	if err != nil {
		return nil, err
	}

	tm := &TwitterManager{
		BaseManager: base,
	}

	// Apply twitter-specific options
	if err := options.ApplyOptions(tm, twitterOpts...); err != nil {
		return nil, err
	}

	return tm, nil
}

// GetID returns the manager's unique identifier
func (tm *TwitterManager) GetID() id.ManagerID {
	return manager.TwitterManagerID
}

// GetDependencies returns the manager IDs that TwitterManager depends on.
// Dependencies include:
// - InsightManagerID: For accessing conversation and user insights
// - PersonalityManagerID: For maintaining consistent personality in responses
func (tm *TwitterManager) GetDependencies() []id.ManagerID {
	return []id.ManagerID{manager.InsightManagerID, manager.PersonalityManagerID}
}

// Process processes incoming tweets and builds conversation context.
// It performs the following steps:
// 1. Verifies the message is from Twitter
// 2. Decodes tweet metadata from the message
// 3. Reconstructs and stores the conversation thread
func (tm *TwitterManager) Process(state *state.State) error {
	tm.Logger.Infof("Executing twitter analysis")

	var metadata twitter.ParsedTweet
	if err := utils.DecodeTweetMetadata(state.Input.Metadata, &metadata); err != nil {
		return fmt.Errorf("failed to unmarshal tweet metadata: %w", err)
	}

	tm.Logger.Infof("Metadata: %v", metadata)

	// Store the tweet thread for context
	if err := tm.storeTweetThread(&metadata); err != nil {
		return fmt.Errorf("failed to store tweet thread: %w", err)
	}

	return nil
}

// PostProcess performs Twitter-specific actions based on the agent's response.
// It handles:
// 1. Verifying the platform is Twitter
// 2. Decoding tweet metadata
// 3. Determining appropriate actions (tweet, like, etc.)
// 4. Executing the decided actions
func (tm *TwitterManager) PostProcess(state *state.State) error {
	tm.Logger.Infof("Executing twitter action")

	response := state.Output

	var parsedTweet twitter.ParsedTweet
	if err := utils.DecodeTweetMetadata(response.Metadata, &parsedTweet); err != nil {
		return fmt.Errorf("failed to unmarshal tweet metadata: %w", err)
	}

	// Get LLM decision on actions
	actions, err := tm.decideTwitterActions(state)
	if err != nil {
		return fmt.Errorf("failed to decide actions: %w", err)
	}

	// Execute the decided actions
	if actions.ShouldTweet {
		if err := tm.sendTweet(response, &parsedTweet); err != nil {
			return fmt.Errorf("failed to send tweet: %w", err)
		}
	}

	if actions.ShouldLike {
		if err := tm.twitterClient.FavoriteTweet(parsedTweet.InReplyToTweetID); err != nil {
			return fmt.Errorf("failed to like tweet: %w", err)
		}
	}

	return nil
}

// Context formats Twitter conversation data for template rendering.
// It performs:
// 1. Decoding tweet metadata from the current message
// 2. Processing and formatting the conversation thread
// 3. Returning formatted conversation data for template use
func (tm *TwitterManager) Context(currentState *state.State) ([]state.StateData, error) {
	var currentTweet twitter.ParsedTweet
	if err := utils.DecodeTweetMetadata(currentState.Input.Metadata, &currentTweet); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tweet metadata: %w", err)
	}

	conversation, err := tm.processAndFormatConversation(
		append(currentState.RecentInteractions, *currentState.Input),
		currentTweet,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to process and format conversation: %w", err)
	}

	return []state.StateData{
		{
			Key:   TwitterConversations,
			Value: conversation,
		},
	}, nil
}

// Store persists a message fragment to storage
func (tm *TwitterManager) Store(fragment *db.Fragment) error {
	// Implement
	return nil
}

// StartBackgroundProcesses initializes any background tasks.
// Note: Currently unimplemented as Twitter interactions are handled reactively
func (tm *TwitterManager) StartBackgroundProcesses() {
	// structurally speaking, we should be fetching tweets here
	// but this goes against how we want to handle twitter
}

// StopBackgroundProcesses cleanly shuts down any background tasks.
// Note: Currently unimplemented as no background processes exist
func (tm *TwitterManager) StopBackgroundProcesses() {
	// Implement
}
