package utils

type Pair[T any] struct {
	First  T
	Second T
}

type Pairs []Pair[any]

func Np[T any](k, s T) Pair[T] {
	return Pair[T]{First: k, Second: s}
}

func (p Pairs) Invalid() {
	f := false
	for _, v := range p {
		if v.First != "" {
			f = true
		} else if f {
			panic("When a parameter name is specified, each subsequent parameter must be specified")
		}
	}
}
