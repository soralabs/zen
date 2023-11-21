package personality

import (
	"zen/internal/managers"
	"zen/internal/state"
	"zen/pkg/options"
)

// Package personality provides functionality for managing agent personalities

// BasePersonality is the key for accessing the core personality configuration
const (
	BasePersonality state.StateDataKey = "base_personality"
)

// MessageExample represents a single example message with its speaker and content
// Used to provide concrete examples of how the personality should communicate
type MessageExample struct {
	User    string // The speaker of the message
	Content string // The actual message content
}

// Personality defines the complete personality configuration for an agent
type Personality struct {
	Name        string // Unique identifier for this personality
	Description string // Detailed description of the personality traits and behavior

	Style      []string // List of communication style characteristics
	Traits     []string // Core personality traits that define behavior
	Background []string // Contextual information about the personality's history
	Expertise  []string // Areas of knowledge and capabilities

	// Example interactions demonstrating personality in action
	ConversationExamples [][]MessageExample // Multi-turn conversation examples
	MessageExamples      []MessageExample   // Single message examples
}

// PersonalityManager handles personality-based behavior and responses
type PersonalityManager struct {
	*managers.BaseManager
	options.RequiredFields

	personality *Personality // Active personality configuration
}
