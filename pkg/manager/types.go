package manager

import (
	"context"

	"github.com/soralabs/zen/pkg/db"
	"github.com/soralabs/zen/pkg/state"

	"github.com/soralabs/zen/pkg/cache"
	"github.com/soralabs/zen/pkg/id"
	"github.com/soralabs/zen/pkg/llm"
	"github.com/soralabs/zen/pkg/logger"
	"github.com/soralabs/zen/pkg/options"
	"github.com/soralabs/zen/pkg/stores"
)

// Package managers provides the core management functionality for agent behaviors
// through a modular system of specialized managers that handle different aspects
// of agent functionality.

type EventType string

// EventData encapsulates event information passed between managers
type EventData struct {
	EventType EventType   // Type of the event being triggered
	Data      interface{} // Associated data payload for the event
}

// EventCallbackFunc defines the signature for event handler functions
type EventCallbackFunc func(eventData EventData) error

// Manager defines the core interface that all managers must implement
// It provides methods for state management, execution, and event handling
type Manager interface {
	// GetID returns the unique identifier for this manager
	GetID() ManagerID

	// GetDependencies returns a list of other manager IDs that this manager depends on
	GetDependencies() []ManagerID

	// Process performs analysis on the current state
	// This is called explicitly by the agent during message processing
	Process(state *state.State) error

	// PostProcess performs actions based on the current state
	// This is called explicitly by the agent after response generation
	PostProcess(state *state.State) error

	// ProvideContext collects and returns manager-specific data for the current state
	Context(state *state.State) ([]state.StateData, error)

	// StoreFragment persists a message fragment to storage
	Store(fragment *db.Fragment) error

	// StartBackgroundProcesses initializes any background tasks
	StartBackgroundProcesses()

	// StopBackgroundProcesses cleanly shuts down any background tasks
	StopBackgroundProcesses()

	// RegisterEventHandler sets up a callback for handling manager events
	RegisterEventHandler(callback EventCallbackFunc)

	// triggerEvent sends an event to the registered handler
	triggerEvent(eventData EventData)
}

// ManagerID is a unique identifier for manager instances
type ManagerID string

// BaseManager provides common functionality for all manager implementations
type BaseManager struct {
	Manager // Embedded interface for type safety
	options.RequiredFields

	AssistantName string
	AssistantID   id.ID

	Ctx           context.Context
	FragmentStore *stores.FragmentStore

	ActorStore   *stores.ActorStore
	SessionStore *stores.SessionStore

	InteractionFragmentStore *stores.FragmentStore

	Cache *cache.Cache

	LLM *llm.LLMClient

	Logger       *logger.Logger
	eventHandler EventCallbackFunc
}
