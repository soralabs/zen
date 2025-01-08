package twitter_manager

import (
	"fmt"
	"strings"
<<<<<<< Updated upstream:managers/twitter/helpers.go
=======
<<<<<<< Updated upstream:internal/managers/twitter/helpers.go
	"zen/internal/db"
	"zen/internal/managers/insight"
	"zen/internal/managers/personality"
	"zen/internal/state"
	"zen/internal/utils"
	"zen/pkg/id"
	"zen/pkg/llm"
	"zen/pkg/twitter"
=======
>>>>>>> Stashed changes:internal/managers/twitter/helpers.go

	"github.com/soralabs/zen/db"
	"github.com/soralabs/zen/id"
	"github.com/soralabs/zen/internal/utils"
<<<<<<< Updated upstream:managers/twitter/helpers.go
	"github.com/soralabs/zen/llm"
	"github.com/soralabs/zen/managers/insight"
	"github.com/soralabs/zen/managers/personality"
	"github.com/soralabs/zen/pkg/twitter"
	"github.com/soralabs/zen/state"
=======
	"github.com/soralabs/zen/pkg/twitter"
	"github.com/soralabs/zen/state"
>>>>>>> Stashed changes:managers/twitter/helpers.go
>>>>>>> Stashed changes:internal/managers/twitter/helpers.go

	"github.com/mitchellh/mapstructure"
	"golang.org/x/sync/errgroup"
)

// storeTweetThread reconstructs and persists a complete conversation thread from Twitter.
// The function performs sequential operations:
// 1. Retrieves the conversation's root tweet
// 2. Fetches all associated replies
// 3. Reconstructs the conversation hierarchy through reply chain traversal
// 4. Persists tweets as message fragments with associated metadata and embeddings
// 5. Manages user entity creation/updates for conversation participants
func (tm *TwitterManager) storeTweetThread(currentTweet *twitter.ParsedTweet) error {
	tm.Logger.WithFields(map[string]interface{}{
		"tweet_id":        currentTweet.TweetID,
		"conversation_id": currentTweet.TweetConversationID,
	}).Infof("Starting to store tweet thread")

	// Track the current branch of the conversation
	var conversationChain []*twitter.ParsedTweet
	conversationID := currentTweet.TweetConversationID

	// Get root tweet details which contains the full conversation
	rootTweet, err := tm.twitterClient.GetTweetDetails(conversationID)
	if err != nil {
		return fmt.Errorf("failed to get root tweet details: %w", err)
	}

	// Parse the root tweet and add it to the chain
	rootParsedTweet, err := tm.twitterClient.ParseTweet(rootTweet)
	if err != nil {
		return fmt.Errorf("failed to parse root tweet: %w", err)
	}
	conversationChain = append(conversationChain, rootParsedTweet)

	// Parse all tweets from the conversation
	replies, err := tm.twitterClient.ParseTweetReplies(rootTweet, "")
	if err != nil {
		return fmt.Errorf("failed to parse tweet replies: %w", err)
	}

	tm.Logger.WithFields(map[string]interface{}{
		"replies": len(replies),
	}).Infof("Found replies")

	// Build map of tweets by ID for easy lookup
	tweetMap := make(map[string]*twitter.ParsedTweet)
	for _, reply := range replies {
		reply := reply // Create new variable for map
		if reply.TweetID == currentTweet.TweetID {
			continue
		}
		tweetMap[reply.TweetID] = &reply
	}

	// Walk up the reply chain from current tweet to root
	currentID := currentTweet.InReplyToTweetID
	for currentID != "" {
		if tweet, ok := tweetMap[currentID]; ok {
			conversationChain = append([]*twitter.ParsedTweet{tweet}, conversationChain...)
			currentID = tweet.InReplyToTweetID
		} else {
			break
		}
	}

	// Format conversation chain for logging
	chainInfo := make([]string, 0, len(conversationChain))
	for _, tweet := range conversationChain {
		chainInfo = append(chainInfo, fmt.Sprintf("@%s: %s", tweet.UserName, tweet.TweetText))
	}

	tm.Logger.WithFields(map[string]interface{}{
		"conversation_chain": strings.Join(chainInfo, " -> "),
	}).Infof("Storing tweet thread")

	// Store tweets in parallel
	var errGroup errgroup.Group
	for _, tweet := range conversationChain {
		tweet := tweet // Create new variable for goroutine
		errGroup.Go(func() error {
			// Determine if this tweet is from the agent
			isAgentTweet := strings.ToLower(tweet.UserName) == strings.ToLower(tm.twitterUsername)
			userID := tm.AssistantID
			if !isAgentTweet {
				userID = id.FromString(tweet.UserID)
				// Store user only if it's not the agent
				if err := tm.ActorStore.Upsert(&db.Actor{
					ID:   userID,
					Name: tweet.UserName,
				}); err != nil {
					return fmt.Errorf("failed to upsert user: %w", err)
				}
			}

			embedding, err := tm.LLM.EmbedText(tweet.TweetText)
			if err != nil {
				return fmt.Errorf("failed to embed tweet text: %w", err)
			}

			// store tweet as a fragment
			tweetFragment, err := utils.CreateTweetFragment(tweet, userID, embedding)
			if err != nil {
				return fmt.Errorf("failed to create tweet fragment: %w", err)
			}

			return tm.InteractionFragmentStore.Upsert(tweetFragment)
		})
	}

	return errGroup.Wait()
}

// processAndFormatConversation structures Twitter conversations for LLM processing.
// The function implements conversation tree analysis through:
// 1. Branch identification and separation
// 2. Chronological message ordering within branches
// 3. Hierarchical conversation reconstruction
// 4. Structured formatting with metadata retention
// 5. Current context demarcation
func (tm *TwitterManager) processAndFormatConversation(recentMessages []db.Fragment, currentTweet twitter.ParsedTweet) (string, error) {
	// Debug print all messages

	tm.Logger.Info("Recent messages:")
	for _, msg := range recentMessages {
		var metadata twitter.ParsedTweet
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			TagName: "json",
			Result:  &metadata,
		})
		if err != nil {
			return "", fmt.Errorf("failed to create decoder: %w", err)
		}
		if err := decoder.Decode(msg.Metadata); err != nil {
			return "", fmt.Errorf("failed to decode tweet metadata: %w", err)
		}
		tm.Logger.WithFields(map[string]interface{}{
			"tweet_id":        metadata.TweetID,
			"in_reply_to":     metadata.InReplyToTweetID,
			"conversation_id": metadata.TweetConversationID,
			"user_name":       metadata.UserName,
			"content":         msg.Content,
			"created_at":      msg.CreatedAt,
		}).Info("Message")
	}

	branches := make(map[string][]db.Fragment)
	currentBranchID := ""

	// First, find all leaf nodes (tweets that aren't replied to)
	leafTweets := make(map[string]db.Fragment)
	repliedToTweets := make(map[string]bool)

	// Build lookup maps
	for _, msg := range recentMessages {
		var metadata twitter.ParsedTweet
		if err := utils.DecodeTweetMetadata(msg.Metadata, &metadata); err != nil {
			return "", fmt.Errorf("failed to decode tweet metadata: %w", err)
		}

		if metadata.InReplyToTweetID != "" {
			repliedToTweets[metadata.InReplyToTweetID] = true
		}

		if metadata.TweetID == currentTweet.TweetID {
			currentBranchID = metadata.TweetID
		}

		leafTweets[metadata.TweetID] = msg
	}

	// Remove tweets that are replied to from leaf nodes
	for tweetID := range repliedToTweets {
		delete(leafTweets, tweetID)
	}

	// Create message lookup for faster access
	msgLookup := make(map[string]db.Fragment)
	for _, msg := range recentMessages {
		var metadata twitter.ParsedTweet
		if err := utils.DecodeTweetMetadata(msg.Metadata, &metadata); err != nil {
			return "", fmt.Errorf("failed to decode tweet metadata: %w", err)
		}
		msgLookup[metadata.TweetID] = msg
	}

	// Build branches from each leaf node
	for leafID, leafMsg := range leafTweets {
		branch := []db.Fragment{leafMsg}
		var metadata twitter.ParsedTweet
		if err := utils.DecodeTweetMetadata(leafMsg.Metadata, &metadata); err != nil {
			return "", fmt.Errorf("failed to decode tweet metadata: %w", err)
		}
		currentID := metadata.InReplyToTweetID

		depth := 0
		maxDepth := 50 // reasonable limit for Twitter threads

		for currentID != "" && depth < maxDepth {
			if msg, exists := msgLookup[currentID]; exists {
				branch = append([]db.Fragment{msg}, branch...)
				var metadata twitter.ParsedTweet
				if err := utils.DecodeTweetMetadata(msg.Metadata, &metadata); err != nil {
					return "", fmt.Errorf("failed to decode tweet metadata: %w", err)
				}
				currentID = metadata.InReplyToTweetID
			} else {
				break
			}
			depth++
		}

		branches[leafID] = branch
	}

	// Format the conversation
	var conversationBuilder strings.Builder

	// Process related branches first
	for branchID, msgs := range branches {
		if branchID != currentBranchID {
			tm.formatBranch(&conversationBuilder, msgs, "RELATED BRANCH", "  ", nil)
		}
	}

	// Process current branch last
	if currentBranch, exists := branches[currentBranchID]; exists {
		tm.formatBranch(&conversationBuilder, currentBranch, "CURRENT BRANCH >>>", "  ", func(msg db.Fragment) string {
			var metadata twitter.ParsedTweet
			if err := utils.DecodeTweetMetadata(msg.Metadata, &metadata); err != nil {
				return ""
			}

			if currentTweet.TweetID == metadata.TweetID {
				return "→ "
			}
			return "  "
		})
	}

	return conversationBuilder.String(), nil
}

// formatBranch applies consistent formatting to a conversation branch.
// Implements customizable message presentation through prefix functions
// for specialized message highlighting and branch organization.
func (tm *TwitterManager) formatBranch(builder *strings.Builder, msgs []db.Fragment, header string, defaultPrefix string, prefixFunc func(db.Fragment) string) {
	builder.WriteString(fmt.Sprintf("\n[%s]\n", header))
	for _, msg := range msgs {
		userName := tm.formatUserName(msg)
		prefix := defaultPrefix
		if prefixFunc != nil {
			prefix = prefixFunc(msg)
		}

		builder.WriteString(fmt.Sprintf("%s[%s] @%s: %s\n",
			prefix,
			msg.CreatedAt.Format("15:04:05"),
			userName,
			msg.Content,
		))
	}
}

// formatUserName standardizes username presentation within conversations.
// Implements agent name normalization and metadata validation.
func (tm *TwitterManager) formatUserName(msg db.Fragment) string {
	var metadata twitter.ParsedTweet
	if err := utils.DecodeTweetMetadata(msg.Metadata, &metadata); err != nil {
		return string(msg.ActorID)
	}
	userName := metadata.UserName

	isAgent := msg.ActorID == tm.AssistantID || strings.ToLower(userName) == strings.ToLower(tm.twitterUsername)
	if isAgent {
		return "YOU"
	}
	return userName
}

// decideTwitterActions determines appropriate platform actions through LLM analysis.
// Evaluates multiple contextual factors:
// 1. Conversation context and derived insights
// 2. Historical user interaction patterns
// 3. Agent personality configuration
// 4. Platform-specific interaction guidelines
//
// The function implements structured decision-making for engagement actions
// including replies and endorsements.
func (tm *TwitterManager) decideTwitterActions(state *state.State) (*TwitterActions, error) {
	// 	sessionInsights, exists := state.GetManagerData(insight.SessionInsights)
	// 	if !exists {
	// 		return nil, fmt.Errorf("session insights not found")
	// 	}
	// 	actorInsights, exists := state.GetManagerData(insight.ActorInsights)
	// 	if !exists {
	// 		return nil, fmt.Errorf("actor insights not found")
	// 	}
	// 	uniqueInsights, exists := state.GetManagerData(insight.UniqueInsights)
	// 	if !exists {
	// 		return nil, fmt.Errorf("unique insights not found")
	// 	}
	// 	basePersonality, exists := state.GetManagerData(personality.BasePersonality)
	// 	if !exists {
	// 		return nil, fmt.Errorf("base personality not found")
	// 	}
	// 	sessionInsights, exists := state.GetManagerData(insight.SessionInsights)
	// 	if !exists {
	// 		return nil, fmt.Errorf("session insights not found")
	// 	}
	// 	actorInsights, exists := state.GetManagerData(insight.ActorInsights)
	// 	if !exists {
	// 		return nil, fmt.Errorf("actor insights not found")
	// 	}
	// 	uniqueInsights, exists := state.GetManagerData(insight.UniqueInsights)
	// 	if !exists {
	// 		return nil, fmt.Errorf("unique insights not found")
	// 	}
	// 	basePersonality, exists := state.GetManagerData(personality.BasePersonality)
	// 	if !exists {
	// 		return nil, fmt.Errorf("base personality not found")
	// 	}

	// 	var actions TwitterActions
	// 	err := tm.LLM.GenerateStructuredOutput(llm.StructuredOutputRequest{
	// 		Messages: []llm.Message{
	// 			llm.NewSystemMessage(`You are an AI that decides which Twitter actions to take based on context and personality. Be thoughtful and deliberate in your decisions.
	// Guidelines:
	// 1. Tweet if a direct response is warranted
	// 2. Like if showing appreciation or acknowledgment is appropriate
	// 3. Can select any combination of the actions, or none at all
	// 	var actions TwitterActions
	// 	err := tm.LLM.GenerateStructuredOutput(llm.StructuredOutputRequest{
	// 		Messages: []llm.Message{
	// 			llm.NewSystemMessage(`You are an AI that decides which Twitter actions to take based on context and personality. Be thoughtful and deliberate in your decisions.
	// Guidelines:
	// 1. Tweet if a direct response is warranted
	// 2. Like if showing appreciation or acknowledgment is appropriate
	// 3. Can select any combination of the actions, or none at all

	// Task:
	// Decide which actions to take based on the above guidelines and the following context.
	// Task:
	// Decide which actions to take based on the above guidelines and the following context.

	// Your personality context:
	// ` + fmt.Sprintf("%v", basePersonality) + `
	// Your personality context:
	// ` + fmt.Sprintf("%v", basePersonality) + `

	// Tweet thread insights (session = conversation):
	// ` + fmt.Sprintf("%v", sessionInsights) + `
	// Tweet thread insights (session = conversation):
	// ` + fmt.Sprintf("%v", sessionInsights) + `

	// User insights (actor = user):
	// ` + fmt.Sprintf("%v", actorInsights) + `
	// User insights (actor = user):
	// ` + fmt.Sprintf("%v", actorInsights) + `

	// Unique insights:
	// ` + fmt.Sprintf("%v", uniqueInsights)),
	// 		},
	// 		ModelType:    llm.ModelTypeAdvanced,
	// 		SchemaName:   "twitter_actions",
	// 		StrictSchema: true,
	// 	}, &actions)
	// Unique insights:
	// ` + fmt.Sprintf("%v", uniqueInsights)),
	// 		},
	// 		ModelType:    llm.ModelTypeAdvanced,
	// 		SchemaName:   "twitter_actions",
	// 		StrictSchema: true,
	// 	}, &actions)

	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to get LLM decision: %w", err)
	// 	}
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to get LLM decision: %w", err)
	// 	}

	// return &actions, nil
	return &TwitterActions{
		ShouldTweet: true,
	}, nil
	// return &actions, nil
	return &TwitterActions{
		ShouldTweet: true,
	}, nil
}

// sendTweet executes Twitter API operations and maintains local state consistency.
// Implements a three-phase operation:
// 1. API interaction for tweet publication
// 2. Metadata synchronization with response data
// 3. Local fragment identifier alignment
func (tm *TwitterManager) sendTweet(response *db.Fragment, parsedTweet *twitter.ParsedTweet) error {
	options := &twitter.TweetOptions{}

	if len(parsedTweet.InReplyToTweetID) > 0 {
		options.ReplyToTweetID = parsedTweet.InReplyToTweetID
	}

	tweet, err := tm.twitterClient.CreateTweet(response.Content, options)
	options := &twitter.TweetOptions{}

	if len(parsedTweet.InReplyToTweetID) > 0 {
		options.ReplyToTweetID = parsedTweet.InReplyToTweetID
	}

	tweet, err := tm.twitterClient.CreateTweet(response.Content, options)
	if err != nil {
		return fmt.Errorf("failed to send tweet: %w", err)
	}

	// Update metadata with tweet ID
	parsedTweet.TweetID = tweet.Data.CreateTweet.TweetResults.Result.RestID

	var metadata db.Metadata
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  &metadata,
	})
	if err != nil {
		return fmt.Errorf("failed to create decoder: %w", err)
	}
	if err := decoder.Decode(parsedTweet); err != nil {
		return fmt.Errorf("failed to decode tweet metadata: %w", err)
	}

	if err := tm.InteractionFragmentStore.UpdateMetadata(response.ID, metadata); err != nil {
		return fmt.Errorf("failed to update tweet metadata: %w", err)
	}

	if err := tm.InteractionFragmentStore.UpdateID(response.ID, id.FromString(parsedTweet.TweetID)); err != nil {
		return fmt.Errorf("failed to update tweet id: %w", err)
	}

	return nil
}
