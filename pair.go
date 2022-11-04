package astool

type Pair[T any] struct {
	First  T
	Second T
}

func MakePair[T any](first, second T) Pair[T] {
	return Pair[T]{first, second}
}
