package personality

import (
	"github.com/soralabs/zen/pkg/db"
	"github.com/soralabs/zen/pkg/manager"
	"github.com/soralabs/zen/pkg/options"
	"github.com/soralabs/zen/pkg/state"
)

// Package personality provides functionality for managing agent personalities and behaviors
// It handles the configuration, storage, and application of personality traits that define
// how an agent interacts in conversations.

// NewPersonalityManager creates a new instance of PersonalityManager with the provided options.
// It initializes the base manager and applies personality-specific configurations.
// baseOpts configure the underlying BaseManager
// personalityOpts configure personality-specific behavior
func NewPersonalityManager(
	baseOpts []options.Option[manager.BaseManager],
	personalityOpts ...options.Option[PersonalityManager],
) (*PersonalityManager, error) {
	base, err := manager.NewBaseManager(baseOpts...)
	if err != nil {
		return nil, err
	}

	pm := &PersonalityManager{
		BaseManager: base,
	}

	// Apply personality-specific options
	if err := options.ApplyOptions(pm, personalityOpts...); err != nil {
		return nil, err
	}

	return pm, nil
}

// GetID returns the manager's unique identifier
func (pm *PersonalityManager) GetID() manager.ManagerID {
	return manager.PersonalityManagerID
}

// GetDependencies returns an empty slice as PersonalityManager has no dependencies
func (pm *PersonalityManager) GetDependencies() []manager.ManagerID {
	return []manager.ManagerID{}
}

// Process analyzes messages for personality-relevant information
// Currently unimplemented as personality is statically configured
func (pm *PersonalityManager) Process(currentState *state.State) error {
	// Implement
	return nil
}

// PostProcess performs personality-driven actions
// Currently unimplemented as personality is applied during response generation
func (pm *PersonalityManager) PostProcess(currentState *state.State) error {
	// Implement
	return nil
}

// Context formats and returns the current personality configuration
// This is used in the prompt template to guide agent behavior
func (pm *PersonalityManager) Context(currentState *state.State) ([]state.StateData, error) {
	// Implement
	personality := formatPersonality(pm.personality)
	return []state.StateData{
		{
			Key:   BasePersonality,
			Value: personality,
		},
	}, nil
}

// Store persists a message fragment to storage
// Currently unimplemented as personality configuration is static
func (pm *PersonalityManager) Store(fragment *db.Fragment) error {
	// Implement
	return nil
}

// StartBackgroundProcesses initializes any background tasks
// Currently unimplemented as personality is statically configured
func (pm *PersonalityManager) StartBackgroundProcesses() {
	// Implement
}

// StopBackgroundProcesses cleanly shuts down any background tasks
// Currently unimplemented as no background processes exist
func (pm *PersonalityManager) StopBackgroundProcesses() {
	// Implement
}
