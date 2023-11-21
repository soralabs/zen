package personality

import (
	"fmt"
	"zen/pkg/options"
)

func (pm *PersonalityManager) ValidateRequiredFields() error {
	if err := pm.BaseManager.ValidateRequiredFields(); err != nil {
		return err
	}
	if pm.personality == nil {
		return fmt.Errorf("personality is required")
	}
	return nil
}

func WithPersonality(personality *Personality) options.Option[PersonalityManager] {
	return func(m *PersonalityManager) error {
		m.personality = personality
		return nil
	}
}
