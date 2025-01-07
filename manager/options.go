package manager

import (
	"context"
	"fmt"

	"github.com/soralabs/zen/id"
	"github.com/soralabs/zen/llm"
	"github.com/soralabs/zen/logger"
	"github.com/soralabs/zen/options"
	"github.com/soralabs/zen/stores"
)

// Package managers provides core functionality for agent behavior management

// ValidateRequiredFields ensures all required fields are set on the BaseManager
// Returns an error if any required field is missing
func (m *BaseManager) ValidateRequiredFields() error {
	if m.Ctx == nil {
		return fmt.Errorf("context is required")
	}
	if m.FragmentStore == nil {
		return fmt.Errorf("fragment store is required")
	}
	if m.ActorStore == nil {
		return fmt.Errorf("actor store is required")
	}
	if m.SessionStore == nil {
		return fmt.Errorf("session store is required")
	}
	if m.LLM == nil {
		return fmt.Errorf("LLM is required")
	}
	if m.Logger == nil {
		return fmt.Errorf("logger is required")
	}
	if m.InteractionFragmentStore == nil {
		return fmt.Errorf("interaction fragment store is required")
	}
	if m.AssistantName == "" {
		return fmt.Errorf("assistant name is required")
	}
	if m.AssistantID == "" {
		return fmt.Errorf("assistant ID is required")
	}
	return nil
}

// WithContext sets the context for the manager
// The context is used for cancellation and timeout control
func WithContext(ctx context.Context) options.Option[BaseManager] {
	return func(m *BaseManager) error {
		m.Ctx = ctx
		return nil
	}
}

// WithCoreDetails sets the core name and ID for the manager
// These are used to identify the assistant in logs and data storage
func WithAssistantDetails(assistantName string, assistantID id.ID) options.Option[BaseManager] {
	return func(m *BaseManager) error {
		m.AssistantName = assistantName
		m.AssistantID = assistantID
		return nil
	}
}

// WithFragmentStore sets the fragment store for the manager
// Used for persisting message fragments
func WithFragmentStore(store *stores.FragmentStore) options.Option[BaseManager] {
	return func(m *BaseManager) error {
		m.FragmentStore = store
		return nil
	}
}

// WithSessionStore sets the session store for the manager
// Used for managing session state and history
func WithSessionStore(store *stores.SessionStore) options.Option[BaseManager] {
	return func(m *BaseManager) error {
		m.SessionStore = store
		return nil
	}
}

// WithActorStore sets the actor store for the manager
// Used for managing actor data and preferences
func WithActorStore(store *stores.ActorStore) options.Option[BaseManager] {
	return func(m *BaseManager) error {
		m.ActorStore = store
		return nil
	}
}

// WithLogger sets the logger instance for the manager
// Used for debugging and monitoring manager operations
func WithLogger(logger *logger.Logger) options.Option[BaseManager] {
	return func(m *BaseManager) error {
		m.Logger = logger
		return nil
	}
}

// WithLLM sets the LLM client for the manager
// Used for generating responses and processing text
func WithLLM(llm *llm.LLMClient) options.Option[BaseManager] {
	return func(m *BaseManager) error {
		m.LLM = llm
		return nil
	}
}

// WithInteractionFragmentStore sets the interaction fragment store for the manager
// Used for storing and retrieving conversation messages
func WithInteractionFragmentStore(store *stores.FragmentStore) options.Option[BaseManager] {
	return func(m *BaseManager) error {
		m.InteractionFragmentStore = store
		return nil
	}
}
