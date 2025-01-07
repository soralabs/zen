package insight

import (
	"fmt"
	"time"

	"github.com/soralabs/zen/cache"
	"github.com/soralabs/zen/db"
	"github.com/soralabs/zen/id"
	"github.com/soralabs/zen/llm"
	"github.com/soralabs/zen/manager"
	"github.com/soralabs/zen/options"
	"github.com/soralabs/zen/state"
	"github.com/soralabs/zen/stores"

	"github.com/pgvector/pgvector-go"
)

// NewInsightManager creates a new instance of InsightManager with the provided options.
// It initializes the base manager and applies insight-specific configurations.
// baseOpts configure the underlying BaseManager
// insightOpts configure insight-specific behavior
func NewInsightManager(
	baseOpts []options.Option[manager.BaseManager],
	insightOpts ...options.Option[InsightManager],
) (*InsightManager, error) {
	base, err := manager.NewBaseManager(baseOpts...)
	if err != nil {
		return nil, err
	}

	im := &InsightManager{
		BaseManager: base,
		examples:    getExampleSets(), // Initialize with training examples
	}

	// Apply insight-specific options
	if err := options.ApplyOptions(im, insightOpts...); err != nil {
		return nil, err
	}

	return im, nil
}

// GetID returns the manager's unique identifier
func (im *InsightManager) GetID() manager.ManagerID {
	return manager.InsightManagerID
}

// GetDependencies returns an empty slice as InsightManager has no dependencies
func (im *InsightManager) GetDependencies() []manager.ManagerID {
	return []manager.ManagerID{}
}

// ExecuteAnalysis performs analysis on the current message to extract new insights.
// It processes the message content, updates existing insights, and maintains insight history.
// The analysis includes:
// 1. Fetching existing insights from cache or storage
// 2. Getting recent messages for context
// 3. Analyzing the session for new insights
// 4. Storing and caching the updated insights
func (im *InsightManager) Process(currentState *state.State) error {
	cacheKey := cache.CacheKey(fmt.Sprintf("%s_%s",
		string(SessionInsights),
		currentState.Input.SessionID,
	))

	// Always fetch fresh data for analysis
	insightData, err := im.getInsightData(currentState.Input)
	if err != nil {
		return err
	}

	// Get recent messages for context
	recentMessages, err := im.InteractionFragmentStore.GetBySession(currentState.Input.SessionID, 20)
	if err != nil {
		return fmt.Errorf("failed to get recent messages: %w", err)
	}

	// Reverse messages for chronological order (oldest first)
	for i, j := 0, len(recentMessages)-1; i < j; i, j = i+1, j-1 {
		recentMessages[i], recentMessages[j] = recentMessages[j], recentMessages[i]
	}

	// Include current message in history
	messageHistory := formatMessageHistory(append(recentMessages, *currentState.Input))

	// Get existing session insights
	sessionInsights, err := im.FragmentStore.SearchByFilter(stores.FragmentFilter{
		SessionID: &currentState.Input.SessionID,
		Metadata: []stores.MetadataCondition{
			{
				Key:      "type",
				Value:    string(SessionInsights),
				Operator: stores.MetadataOpEquals,
			},
		},
		Limit: 10,
	})
	if err != nil {
		return fmt.Errorf("failed to get session insights: %w", err)
	}

	// Get existing actor insights
	actorInsights, err := im.FragmentStore.SearchByFilter(stores.FragmentFilter{
		ActorID: &currentState.Input.ActorID,
		Metadata: []stores.MetadataCondition{
			{
				Key:      "type",
				Value:    string(ActorInsights),
				Operator: stores.MetadataOpEquals,
			},
		},
		Limit: 10,
	})
	if err != nil {
		return fmt.Errorf("failed to get actor insights: %w", err)
	}

	// Format insights for LLM processing
	existingInsights := im.formatAllInsights(sessionInsights, actorInsights)
	exampleInsights := formatExampleInsights(im.examples)

	// Prepare LLM messages for insight extraction
	messages := []llm.Message{
		llm.NewSystemMessage(`You are an insight extraction system that identifies concrete, factual information from conversations.
Focus on extracting clear, verifiable information while maintaining high precision.
Avoid speculation and uncertain information.

IMPORTANT RULES:
1. Do NOT extract insights that are already present in the existing insights list - only extract new information.
2. If you find any existing insights that are contradicted by new information, include their IDs in outdated_insight_ids so they can be removed.
3. Mark insights as outdated if they are no longer accurate or have been superseded by new information.

INSIGHT TYPES:
1. Session insights: Facts specific to this session
2. Actor insights: Facts about the actor that are relevant across all sessions

CATEGORIES:
- preference: Actor preferences and likes/dislikes
- technical_detail: Technical information about projects or code
- personal_info: Facts about the actor themselves
- requirement: Specific needs or requirements
- date_time: Temporal information
- entity: Names of tools, frameworks, or other entities

GUIDELINES:
1. Extract only explicitly stated information
2. Make each insight atomic (one piece of information)
3. Include confidence levels based on clarity of statements
4. Avoid any speculation or assumptions
5. Focus on factual, verifiable information
6. Do not extract insights, if not necessary
7. Actor IDs must be included in all insights, and must be a valid actor ID from the current session`),

		llm.NewUserMessage(fmt.Sprintf("Here are examples of good insights:\n%s",
			exampleInsights)),

		llm.NewUserMessage(fmt.Sprintf("Please analyze this conversation:\n%s\n\nExisting insights:\n%s",
			messageHistory,
			existingInsights)),
	}

	var result InsightResponse
	// Generate insights using LLM
	err = im.LLM.GenerateStructuredOutput(llm.StructuredOutputRequest{
		Messages:     messages,
		ModelType:    llm.ModelTypeAdvanced,
		SchemaName:   "insight_extraction",
		StrictSchema: true,
	}, &result)
	if err != nil {
		return fmt.Errorf("failed to generate insights: %w", err)
	}

	// Store new insights and update embeddings
	for _, insight := range result.NewInsights {
		insightFragment := &db.Fragment{
			ID:        id.New(),
			ActorID:   id.ID(insight.ActorID),
			SessionID: currentState.Input.SessionID,
			Content:   insight.Content,
			Metadata: map[string]interface{}{
				"type":           insight.Type,
				"category":       insight.Category,
				"confidence":     insight.Confidence,
				"source_context": insight.SourceContext,
				"timestamp":      time.Now().Unix(),
			},
		}

		// Generate embedding for semantic search
		embedding, err := im.LLM.EmbedText(insight.Content)
		if err != nil {
			im.Logger.Warnf("Failed to generate embedding for insight: %v", err)
			continue
		}
		insightFragment.Embedding = pgvector.NewVector(embedding)

		if err := im.FragmentStore.Create(insightFragment); err != nil {
			return fmt.Errorf("failed to store new insight: %w", err)
		}

		// Update cached data with new insight
		if insight.Type == string(SessionInsights) {
			insightData.SessionInsights = append(insightData.SessionInsights, *insightFragment)
		} else if insight.Type == string(ActorInsights) {
			insightData.ActorInsights = append(insightData.ActorInsights, *insightFragment)
		}
	}

	// Remove outdated insights
	for _, outdatedID := range result.OutdatedInsightIDs {
		if err := im.FragmentStore.DeleteByID(outdatedID); err != nil {
			im.Logger.Warnf("Failed to delete outdated insight: %v", err)
		}

		// Update cache
		insightData.SessionInsights = removeInsightByID(insightData.SessionInsights, outdatedID)
		insightData.ActorInsights = removeInsightByID(insightData.ActorInsights, outdatedID)
	}

	// Update cache with the modified insight data
	im.Cache.Set(cacheKey, insightData)

	return nil
}

// PostProcess performs post-processing actions
func (im *InsightManager) PostProcess(currentState *state.State) error {
	return nil
}

// ProvideContext collects and formats all relevant insights for the current session
// It retrieves session, user, and similar insights, then formats them for template use
func (im *InsightManager) Context(currentState *state.State) ([]state.StateData, error) {
	cacheKey := cache.CacheKey(fmt.Sprintf("%s_%s",
		string(SessionInsights),
		currentState.Input.SessionID,
	))

	var insightData *cachedInsightData
	if cached, found := im.Cache.Get(cacheKey); found {
		insightData = cached.(*cachedInsightData)
	} else {
		// If not in cache, do a fresh fetch
		var err error
		insightData, err = im.getInsightData(currentState.Input)
		if err != nil {
			return nil, err
		}
	}

	uniqueInsights := getUniqueInsights(insightData.SimilarInsights, insightData.SessionInsights, insightData.ActorInsights)

	return []state.StateData{
		{
			Key:   SessionInsights,
			Value: im.formatInsightFragments(insightData.SessionInsights),
		},
		{
			Key:   ActorInsights,
			Value: im.formatInsightFragments(insightData.ActorInsights),
		},
		{
			Key:   UniqueInsights,
			Value: im.formatInsightFragments(uniqueInsights),
		},
	}, nil
}

// Store persists a new fragment to storage (currently unimplemented)
func (im *InsightManager) Store(fragment *db.Fragment) error {
	return nil
}

// StartBackgroundProcesses initializes any background tasks (currently unimplemented)
func (im *InsightManager) StartBackgroundProcesses() {
	// Implement
}

func (im *InsightManager) StopBackgroundProcesses() {
	// Implement
}
