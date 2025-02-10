package engine

import (
	"fmt"
	"time"

	"github.com/soralabs/zen/db"
	"github.com/soralabs/zen/manager"
	"github.com/soralabs/zen/state"
	"golang.org/x/sync/errgroup"
)

// ProcessBuilder provides a fluent interface for configuring process operations
type ProcessBuilder struct {
	engine        *Engine
	state         *state.State
	managerFilter []manager.ManagerID
	shouldStore   bool
	metadata      map[string]interface{}
	createdAt     *time.Time
	managerOrder  []manager.ManagerID
	validators    []func(*state.State) error
}

// PostProcessBuilder provides a fluent interface for configuring post-process operations
type PostProcessBuilder struct {
	engine        *Engine
	response      *db.Fragment
	state         *state.State
	managerFilter []manager.ManagerID
	shouldStore   bool
	metadata      map[string]interface{}
	createdAt     *time.Time
	managerOrder  []manager.ManagerID
	validators    []func(*state.State) error
}

// NewProcessBuilder starts building a process operation
func (e *Engine) NewProcessBuilder() *ProcessBuilder {
	return &ProcessBuilder{
		engine:      e,
		shouldStore: true, // Default to storing
		metadata:    make(map[string]interface{}),
	}
}

// NewPostProcessBuilder starts building a post-process operation
func (e *Engine) NewPostProcessBuilder() *PostProcessBuilder {
	return &PostProcessBuilder{
		engine:      e,
		shouldStore: true, // Default to storing
		metadata:    make(map[string]interface{}),
	}
}

// WithState sets the state for the process operation
func (b *ProcessBuilder) WithState(state *state.State) *ProcessBuilder {
	b.state = state
	return b
}

// WithManagerFilter sets which managers to execute
func (b *ProcessBuilder) WithManagerFilter(managerFilter []manager.ManagerID) *ProcessBuilder {
	b.managerFilter = managerFilter
	return b
}

// WithManagerOrder sets a custom execution order for managers
func (b *ProcessBuilder) WithManagerOrder(order []manager.ManagerID) *ProcessBuilder {
	b.managerOrder = order
	return b
}

// WithMetadata adds metadata to the input fragment
func (b *ProcessBuilder) WithMetadata(key string, value interface{}) *ProcessBuilder {
	b.metadata[key] = value
	return b
}

// WithCreatedAt sets a custom creation time for the input fragment
func (b *ProcessBuilder) WithCreatedAt(t time.Time) *ProcessBuilder {
	b.createdAt = &t
	return b
}

// WithValidator adds a validation function to be run before processing
func (b *ProcessBuilder) WithValidator(validator func(*state.State) error) *ProcessBuilder {
	b.validators = append(b.validators, validator)
	return b
}

// ShouldStore controls whether to store the input fragment
func (b *ProcessBuilder) ShouldStore(shouldStore bool) *ProcessBuilder {
	b.shouldStore = shouldStore
	return b
}

// WithDefaults sets up common default configurations for process operations:
// - Uses all available managers
// - Uses engine's default manager order
// - Stores the fragment
// - Sets creation time to now
// - Adds default metadata (timestamp, version)
func (b *ProcessBuilder) WithDefaults() *ProcessBuilder {
	// Use all available managers
	var allManagers []manager.ManagerID
	for _, m := range b.engine.managers {
		allManagers = append(allManagers, m.GetID())
	}
	b.managerFilter = allManagers

	// Use engine's manager order if available
	if len(b.engine.managerOrder) > 0 {
		b.managerOrder = make([]manager.ManagerID, len(b.engine.managerOrder))
		copy(b.managerOrder, b.engine.managerOrder)
	}

	// Enable storage
	b.shouldStore = true

	// Set creation time to now
	now := time.Now()
	b.createdAt = &now

	return b
}

// Execute runs the process operation with the configured options
func (b *ProcessBuilder) Execute() error {
	if b.state == nil {
		return fmt.Errorf("state is required")
	}

	// Run validators
	for _, validator := range b.validators {
		if err := validator(b.state); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
	}

	input := b.state.Input
	b.engine.logger.WithFields(map[string]interface{}{
		"input":         input.ID,
		"managerFilter": b.managerFilter,
	}).Info("Processing input with manager filter")

	actor, err := b.engine.actorStore.GetByID(input.ActorID)
	if err != nil {
		return fmt.Errorf("failed to get actor: %w", err)
	}

	session, err := b.engine.sessionStore.GetByID(input.SessionID)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	// Create a copy of the input fragment for storage
	inputCopy := &db.Fragment{
		ID:        input.ID,
		ActorID:   input.ActorID,
		SessionID: input.SessionID,
		Content:   input.Content,
		Metadata:  input.Metadata,
		Embedding: input.Embedding,
		Actor:     actor,
		Session:   session,
	}

	// Apply custom metadata
	if len(b.metadata) > 0 {
		if inputCopy.Metadata == nil {
			inputCopy.Metadata = make(map[string]interface{})
		}
		for k, v := range b.metadata {
			inputCopy.Metadata[k] = v
		}
	}

	// Apply custom creation time or use current time
	if b.createdAt != nil {
		inputCopy.CreatedAt = *b.createdAt
	} else if inputCopy.CreatedAt.IsZero() {
		inputCopy.CreatedAt = time.Now()
	}

	b.state.Input = inputCopy

	// Use custom manager order if specified, otherwise use engine's order
	var executionOrder []manager.ManagerID
	if len(b.managerOrder) > 0 {
		executionOrder = b.managerOrder
	} else if len(b.engine.managerOrder) > 0 {
		executionOrder = b.engine.managerOrder
	} else {
		for _, m := range b.engine.managers {
			if b.engine.containsManager(b.managerFilter, m.GetID()) {
				executionOrder = append(executionOrder, m.GetID())
			}
		}
	}

	managerMap := make(map[manager.ManagerID]manager.Manager)
	for _, m := range b.engine.managers {
		managerMap[m.GetID()] = m
	}

	var errGroup errgroup.Group
	for _, mid := range executionOrder {
		if m, exists := managerMap[mid]; exists {
			mCopy := m // capture variable for closure
			errGroup.Go(func() error {
				return mCopy.Process(b.state)
			})
		}
	}

	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("failed to execute manager processes: %w", err)
	}

	if b.shouldStore {
		if err := b.engine.interactionFragmentStore.Upsert(inputCopy); err != nil {
			return fmt.Errorf("failed to store input: %w", err)
		}
	}

	return nil
}

// WithState sets the state for the post-process operation
func (b *PostProcessBuilder) WithState(state *state.State) *PostProcessBuilder {
	b.state = state
	return b
}

// WithResponse sets the response fragment for post-processing
func (b *PostProcessBuilder) WithResponse(response *db.Fragment) *PostProcessBuilder {
	b.response = response
	return b
}

// WithManagerFilter sets which managers to execute
func (b *PostProcessBuilder) WithManagerFilter(managerFilter []manager.ManagerID) *PostProcessBuilder {
	b.managerFilter = managerFilter
	return b
}

// WithManagerOrder sets a custom execution order for managers
func (b *PostProcessBuilder) WithManagerOrder(order []manager.ManagerID) *PostProcessBuilder {
	b.managerOrder = order
	return b
}

// WithMetadata adds metadata to the response fragment
func (b *PostProcessBuilder) WithMetadata(key string, value interface{}) *PostProcessBuilder {
	b.metadata[key] = value
	return b
}

// WithCreatedAt sets a custom creation time for the response fragment
func (b *PostProcessBuilder) WithCreatedAt(t time.Time) *PostProcessBuilder {
	b.createdAt = &t
	return b
}

// WithValidator adds a validation function to be run before post-processing
func (b *PostProcessBuilder) WithValidator(validator func(*state.State) error) *PostProcessBuilder {
	b.validators = append(b.validators, validator)
	return b
}

// ShouldStore controls whether to store the response fragment
func (b *PostProcessBuilder) ShouldStore(shouldStore bool) *PostProcessBuilder {
	b.shouldStore = shouldStore
	return b
}

// WithDefaults sets up common default configurations for post-process operations:
// - Uses all available managers
// - Uses engine's default manager order
// - Stores the fragment
// - Sets creation time to now
// - Adds default metadata (timestamp, version)
func (b *PostProcessBuilder) WithDefaults() *PostProcessBuilder {
	// Use all available managers
	var allManagers []manager.ManagerID
	for _, m := range b.engine.managers {
		allManagers = append(allManagers, m.GetID())
	}
	b.managerFilter = allManagers

	// Use engine's manager order if available
	if len(b.engine.managerOrder) > 0 {
		b.managerOrder = make([]manager.ManagerID, len(b.engine.managerOrder))
		copy(b.managerOrder, b.engine.managerOrder)
	}

	// Enable storage
	b.shouldStore = true

	// Set creation time to now
	now := time.Now()
	b.createdAt = &now

	return b
}

// Execute runs the post-process operation with the configured options
func (b *PostProcessBuilder) Execute() error {
	if b.state == nil {
		return fmt.Errorf("state is required")
	}
	if b.response == nil {
		return fmt.Errorf("response is required")
	}

	// Run validators
	for _, validator := range b.validators {
		if err := validator(b.state); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
	}

	actor, err := b.engine.actorStore.GetByID(b.response.ActorID)
	if err != nil {
		return fmt.Errorf("failed to get actor: %w", err)
	}

	session, err := b.engine.sessionStore.GetByID(b.response.SessionID)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	responseCopy := &db.Fragment{
		ID:        b.response.ID,
		ActorID:   b.response.ActorID,
		SessionID: b.response.SessionID,
		Content:   b.response.Content,
		Metadata:  b.response.Metadata,
		Embedding: b.response.Embedding,
		Actor:     actor,
		Session:   session,
	}

	// Apply custom metadata
	if len(b.metadata) > 0 {
		if responseCopy.Metadata == nil {
			responseCopy.Metadata = make(map[string]interface{})
		}
		for k, v := range b.metadata {
			responseCopy.Metadata[k] = v
		}
	}

	// Apply custom creation time or use current time
	if b.createdAt != nil {
		responseCopy.CreatedAt = *b.createdAt
	} else if responseCopy.CreatedAt.IsZero() {
		responseCopy.CreatedAt = time.Now()
	}

	b.state.Output = responseCopy

	// Use custom manager order if specified, otherwise use engine's order
	var executionOrder []manager.ManagerID
	if len(b.managerOrder) > 0 {
		executionOrder = b.managerOrder
	} else if len(b.engine.managerOrder) > 0 {
		executionOrder = b.engine.managerOrder
	} else {
		for _, m := range b.engine.managers {
			if b.engine.containsManager(b.managerFilter, m.GetID()) {
				executionOrder = append(executionOrder, m.GetID())
			}
		}
	}

	managerMap := make(map[manager.ManagerID]manager.Manager)
	for _, m := range b.engine.managers {
		managerMap[m.GetID()] = m
	}

	for _, mid := range executionOrder {
		if m, exists := managerMap[mid]; exists {
			if err := m.PostProcess(b.state); err != nil {
				return fmt.Errorf("manager %s failed: %w", mid, err)
			}
		}
	}

	if b.shouldStore {
		if err := b.engine.interactionFragmentStore.Upsert(b.response); err != nil {
			return fmt.Errorf("failed to store response: %w", err)
		}
	}

	return nil
}
