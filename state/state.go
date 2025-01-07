package state

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
