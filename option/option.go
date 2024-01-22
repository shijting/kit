package option

type Option[T any] func(t *T)

type Options[T any] []Option[T]

func (opts Options[T]) Apply(t *T) {
	for _, opt := range opts {
		opt(t)
	}
}

type OptionErr[T any] func(t *T) error

type OptionsErr[T any] []OptionErr[T]

func (opts OptionsErr[T]) Apply(t *T) error {
	for _, opt := range opts {
		if err := opt(t); err != nil {
			return err
		}
	}
	return nil
}
