package managers

import (
	"fmt"
	"time"
	"zen/internal/db"
	"zen/internal/state"

	"zen/pkg/cache"
	"zen/pkg/options"
)

// GetID returns the base manager identifier
func (bm *BaseManager) GetID() ManagerID {
	return BaseManagerID
}

// GetDependencies returns an empty dependency list for the base manager
func (bm *BaseManager) GetDependencies() []ManagerID {
	return []ManagerID{}
}

// Process provides a default implementation that panics
// Managers should override this method with their specific analysis logic
func (bm *BaseManager) Process(state *state.State) error {
	panic("Process not implemented")
}

// PostProcess provides a default implementation that panics
// Managers should override this method with their specific post-processing logic
func (bm *BaseManager) PostProcess(state *state.State) error {
	panic("PostProcess not implemented")
}

// Context provides a default implementation that panics
// Managers should override this method to provide their specific context data
func (bm *BaseManager) Context(state *state.State) ([]state.StateData, error) {
	panic("Context not implemented")
}

// Store persists a fragment to the fragment store
func (bm *BaseManager) Store(fragment *db.Fragment) error {
	return bm.FragmentStore.Create(fragment)
}

// StartBackgroundProcesses provides a default implementation that panics
// Managers should override this method if they need background processing
func (bm *BaseManager) StartBackgroundProcesses() {
	panic("StartBackgroundProcesses not implemented")
}

// StopBackgroundProcesses provides a default implementation that panics
// Managers should override this method if they need to clean up background processes
func (bm *BaseManager) StopBackgroundProcesses() {
	panic("StopBackgroundProcesses not implemented")
}

// RegisterEventHandler sets the event handler callback for this manager
func (bm *BaseManager) RegisterEventHandler(callback EventCallbackFunc) {
	bm.eventHandler = callback
}

// triggerEvent sends an event to the registered handler
// Panics if no handler is registered
func (bm *BaseManager) triggerEvent(eventData EventData) {
	if bm.eventHandler != nil {
		bm.eventHandler(eventData)
	} else {
		panic("No event handler registered")
	}
}

// NewBaseManager creates a new BaseManager instance with the provided options
func NewBaseManager(opts ...options.Option[BaseManager]) (*BaseManager, error) {
	bm := &BaseManager{
		Cache: cache.New(cache.Config{
			MaxSize:       1000,
			TTL:           15 * time.Minute,
			CleanupPeriod: 30 * time.Minute,
		}),
	}
	if err := options.ApplyOptions(bm, opts...); err != nil {
		return nil, fmt.Errorf("failed to create base manager: %w", err)
	}
	return bm, nil
}
