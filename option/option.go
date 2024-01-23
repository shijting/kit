package option

type Option[T any] func(t *T)

type Options[T any] []Option[T]

// Apply applies all options to t.
func (opts Options[T]) Apply(t *T) {
	for _, opt := range opts {
		opt(t)
	}
}

type OptionErr[T any] func(t *T) error

type OptionsErr[T any] []OptionErr[T]

// Apply applies all options to t.
// Returns an error if any of the options returns an error.
// Returns nil if no error occurred.
func (opts OptionsErr[T]) Apply(t *T) error {
	for _, opt := range opts {
		if err := opt(t); err != nil {
			return err
		}
	}
	return nil
}
