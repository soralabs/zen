package state

import "github.com/soralabs/zen/id"

// Package state provides core functionality for managing conversation state and context
// in the agent system. It handles both structured manager data and custom runtime data,
// while providing methods for state manipulation and template-based prompt generation.

// AddManagerData adds a slice of StateData entries to the state's manager data store.
// If the manager data map hasn't been initialized, it creates a new one.
func (s *State) AddManagerData(data []StateData) *State {
	if s.managerData == nil {
		s.managerData = make(map[StateDataKey]interface{})
	}

	for _, d := range data {
		s.managerData[d.Key] = d.Value
	}

	return s
}

// GetManagerData retrieves manager-specific data by its key.
// Returns the value and a boolean indicating if the key exists.
func (s *State) GetManagerData(key StateDataKey) (interface{}, bool) {
	value, exists := s.managerData[key]
	return value, exists
}

// AddCustomData adds a custom key-value pair to the state's custom data store.
// This is useful for platform-specific or temporary data that doesn't fit into manager data.
func (s *State) AddCustomData(key string, value interface{}) *State {
	if s.customData == nil {
		s.customData = make(map[string]interface{})
	}
	s.customData[key] = value

	return s
}

// GetCustomData retrieves a custom data value by its key.
// Returns the value and a boolean indicating if the key exists.
func (s *State) GetCustomData(key string) (interface{}, bool) {
	if s.customData == nil {
		return nil, false
	}
	value, exists := s.customData[key]
	return value, exists
}

// Reset clears all manager and custom data from the state.
// This is typically called before updating the state with fresh data.
func (s *State) Reset() {
	s.managerData = make(map[StateDataKey]interface{})
	s.customData = make(map[string]interface{})
}

// SetManagersToProcess specifies which managers should be executed during Process
func (s *State) SetManagersToProcess(managerIDs ...id.ManagerID) *State {
	s.managerExecution.processManagers = make(map[id.ManagerID]bool)
	for _, id := range managerIDs {
		s.managerExecution.processManagers[id] = true
	}
	return s
}

// SetManagersToPostProcess specifies which managers should be executed during PostProcess
func (s *State) SetManagersToPostProcess(managerIDs ...id.ManagerID) *State {
	s.managerExecution.postProcessManagers = make(map[id.ManagerID]bool)
	for _, id := range managerIDs {
		s.managerExecution.postProcessManagers[id] = true
	}
	return s
}

// ShouldProcessManager checks if a manager should be executed during Process
func (s *State) ShouldProcessManager(managerID id.ManagerID) bool {
	// If no managers are specified, execute all
	if len(s.managerExecution.processManagers) == 0 {
		return true
	}
	return s.managerExecution.processManagers[managerID]
}

// ShouldPostProcessManager checks if a manager should be executed during PostProcess
func (s *State) ShouldPostProcessManager(managerID id.ManagerID) bool {
	// If no managers are specified, execute all
	if len(s.managerExecution.postProcessManagers) == 0 {
		return true
	}
	return s.managerExecution.postProcessManagers[managerID]
}
