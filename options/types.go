package options

// Option is a function that configures a value T
type Option[T any] func(*T) error

// RequiredFields is an interface that ensures a type can validate its required fields
type RequiredFields interface {
	ValidateRequiredFields() error
}
