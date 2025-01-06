package state

import (
	"html/template"
	"zen/internal/db"
	"zen/pkg/llm"

	toolkit "github.com/soralabs/toolkit/go"
)

// Package state provides core state management functionality for the agent system

// StateDataKey represents a unique identifier for state data entries
type StateDataKey string

// StateData represents a key-value pair of data provided by managers
type StateData struct {
	Key   StateDataKey
	Value interface{}
}

// State represents the current context and state of a conversation
// It maintains core conversation data, user information, and both manager and custom data
type State struct {
	// Core conversation data
	Input  *db.Fragment // The current input
	Output *db.Fragment // The LLM response

	// Actor information
	Actor *db.Actor // Information about where it came from

	// Recent data
	RecentInteractions   []db.Fragment
	RelevantInteractions []db.Fragment
	Tools                []toolkit.Tool
	// Manager-specific data storage
	// Stores data provided by various managers keyed by StateDataKey
	managerData map[StateDataKey]interface{}

	// Custom data storage for arbitrary key-value pairs
	// Used for platform-specific or temporary data storage
	customData map[string]interface{}
}

// NewState creates and initializes a new State instance with empty data stores
func NewState() *State {
	return &State{
		managerData: make(map[StateDataKey]interface{}),
		customData:  make(map[string]interface{}),
	}
}

// PromptSection represents a single section of a prompt template with its role and content
type PromptSection struct {
	Role     llm.Role // The role of this section (system, user, assistant, etc)
	Template string   // The template text for this section
	Name     string   // Optional name for the role (e.g., specific user identifiers)
}

// PromptBuilder facilitates the construction of structured prompts
// It manages template sections and associated state data
type PromptBuilder struct {
	state     *State                       // Reference to the current state
	sections  []PromptSection              // Ordered list of prompt sections
	stateData map[StateDataKey]interface{} // Manager-provided data for template rendering
	helpers   template.FuncMap             // Function map for custom template functions
	err       error                        // Tracks any errors during building
}
