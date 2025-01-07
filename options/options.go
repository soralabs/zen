package options

import (
	"fmt"
)

// ApplyOptions applies a list of options to a value
func ApplyOptions[T any](value *T, opts ...Option[T]) error {
	// Apply all options
	for _, opt := range opts {
		if err := opt(value); err != nil {
			return fmt.Errorf("failed to apply option: %w", err)
		}
	}

	// If the type implements RequiredFields, validate them
	if v, ok := any(value).(RequiredFields); ok {
		if err := v.ValidateRequiredFields(); err != nil {
			return fmt.Errorf("missing required fields: %w", err)
		}
	}

	return nil
}

// WithValidation adds validation to any option
func WithValidation[T any](opt Option[T], validationFn func(*T) error) Option[T] {
	return func(value *T) error {
		if err := opt(value); err != nil {
			return err
		}
		return validationFn(value)
	}
}

// WithDefault sets a default value if the option returns an error
func WithDefault[T any, V any](opt Option[T], field *V, defaultValue V) Option[T] {
	return func(value *T) error {
		if err := opt(value); err != nil {
			*field = defaultValue
			return nil
		}
		return nil
	}
}
