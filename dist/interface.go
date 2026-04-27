package dist

type Pool[T any] interface {
	Next() T
}
