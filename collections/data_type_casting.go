package collections

type IntoSlice[T any] interface {
	Into() []T
}
