package option

type Option[T any] func(t *T)

type Options[T any] []Option[T]

func (opts Options[T]) Apply(t *T) {
	for _, opt := range opts {
		opt(t)
	}
}
