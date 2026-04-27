package collections

// IntoSlice is implemented by collections that can materialize their elements into a slice.
type IntoSlice[T any] interface {
	Into() []T
}
