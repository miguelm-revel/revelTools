package dist

// Pool is implemented by any type that distributes items from a collection.
type Pool[T any] interface {
	Next() T
}
