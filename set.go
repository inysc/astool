package astool

import "fmt"

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](keys ...T) Set[T] {
	st := Set[T]{}
	for _, v := range keys {
		st[v] = struct{}{}
	}
	return st
}

func (st Set[T]) Keys() []T {
	keys := make([]T, 0, len(st))
	for k := range st {
		keys = append(keys, k)
	}
	return keys
}

func (st Set[T]) Union(sets ...Set[T]) Set[T] {
	ret := NewSet(st.Keys()...)
	for _, set := range sets {
		for _, key := range set.Keys() {
			ret.Insert(key)
		}
	}
	return ret
}

func (st Set[T]) Intersect(sets ...Set[T]) Set[T] {
	ret := NewSet(st.Keys()...)
	for _, set := range sets {
		for _, key := range ret.Keys() {
			if !set.Has(key) {
				ret.Delete(key)
			}
		}
	}
	return ret
}

func (st Set[T]) Difference(sets ...Set[T]) Set[T] {
	ret := NewSet(st.Keys()...)
	for _, set := range sets {
		for _, key := range set.Keys() {
			ret.Delete(key)
		}
	}
	return ret
}

func (st Set[T]) Equal(set Set[T]) bool {
	if st.Len() != set.Len() {
		return false
	}
	for _, key := range st.Keys() {
		if !set.Has(key) {
			return false
		}
	}
	return true
}

func (st Set[T]) IsSubset(set Set[T]) bool {
	if st.Len() > set.Len() {
		return false
	}
	for _, key := range st.Keys() {
		if !set.Has(key) {
			return false
		}
	}
	return true
}

func (st Set[T]) IsSuperset(set Set[T]) bool {
	return set.IsSubset(st)
}

func (st Set[T]) String() string {
	return fmt.Sprint(st.Keys())
}

func (st Set[T]) Copy() Set[T] {
	return NewSet(st.Keys()...)
}

func (st Set[T]) Clear() {
	st = make(map[T]struct{})
}

func (st Set[T]) Pop() T {
	for key := range st {
		st.Delete(key)
		return key
	}
	panic("pop from an empty set")
}

func (st Set[T]) Insert(key T) {
	(st)[key] = struct{}{}
}

func (st Set[T]) Delete(key T) {
	delete(st, key)
}

func (st Set[T]) Has(key T) bool {
	_, ok := (st)[key]
	return ok
}

func (st Set[T]) Len() int {
	return len(st)
}
