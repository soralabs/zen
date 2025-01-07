package personality

import (
	"fmt"

	"github.com/soralabs/zen/options"
)

// ValidateRequiredFields checks if all required fields, including the personality, are set
func (pm *PersonalityManager) ValidateRequiredFields() error {
	if err := pm.BaseManager.ValidateRequiredFields(); err != nil {
		return err
	}
	if pm.personality == nil {
		return fmt.Errorf("personality is required")
	}
	return nil
}

// WithPersonality returns an Option that sets the Personality on PersonalityManager
func WithPersonality(personality *Personality) options.Option[PersonalityManager] {
	return func(m *PersonalityManager) error {
		m.personality = personality
		return nil
	}
}
