package insight

import (
	"zen/internal/db"
	"zen/internal/managers"
	"zen/internal/state"
	"zen/pkg/id"
	"zen/pkg/options"
)

// Package insight provides functionality for generating and managing conversation insights

// StateDataKey constants for different types of insights
const (
	// ActorInsights represents insights about specific actors
	ActorInsights state.StateDataKey = "actor_insights"

	// SessionInsights represents insights about the current session
	SessionInsights state.StateDataKey = "session_insights"

	// UniqueInsights represents filtered unique insights across all types
	UniqueInsights state.StateDataKey = "unique_insights"
)

// InsightManager handles the generation and management of conversation insights
type InsightManager struct {
	*managers.BaseManager
	options.RequiredFields

	examples []ExampleSet // Training examples for insight generation
}

// ExampleInsight represents a single insight with metadata
type ExampleInsight struct {
	ID            id.ID   // Unique identifier for the insight
	ActorID       id.ID   // ID of the actor this insight relates to
	Content       string  // The actual insight text
	Category      string  // Classification category for the insight
	Confidence    float64 // Confidence score for the insight (0-1)
	SourceContext string  // Original context that generated this insight
	Type          string  // Type of insight (conversation or user)
}

// ExampleSet represents a collection of example messages and their associated insights
type ExampleSet struct {
	Messages []ExampleMessage // Example conversation messages
	Insights []ExampleInsight // Associated insights for the messages
}

type ExampleMessage struct {
	Role    string
	Content string
	TimeAgo string
}

type InsightResponse struct {
	NewInsights []struct {
		Content       string  `json:"content" jsonschema:"required"`
		Category      string  `json:"category" jsonschema:"required,enum=preference,enum=technical_detail,enum=personal_info,enum=requirement,enum=date_time,enum=entity" description:"preference: User preferences and likes/dislikes, technical_detail: Technical specifications or requirements, personal_info: Personal information about the user, requirement: Project or task requirements, date_time: Time-related information, entity: Named entities or specific objects"`
		Confidence    float64 `json:"confidence" jsonschema:"required,minimum=0,maximum=1"`
		SourceContext string  `json:"source_context" jsonschema:"required"`
		ActorID       string  `json:"actor_id" jsonschema:"required"`
		Type          string  `json:"type" jsonschema:"required,enum=session_insights,enum=actor_insights" description:"session_insights: Insight derived from the current session, actor_insights: Persistent insight about the actor"`
	} `json:"new_insights" jsonschema:"required"`
	OutdatedInsightIDs []id.ID `json:"outdated_insight_ids" jsonschema:"required"`
}

type cachedInsightData struct {
	SessionInsights []db.Fragment
	ActorInsights   []db.Fragment
	SimilarInsights []db.Fragment
}
