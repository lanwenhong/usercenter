package util

import "context"

type Set[T comparable] struct {
	InMap map[T]interface{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		InMap: map[T]interface{}{},
	}
}

func (s *Set[T]) Set(ctx context.Context, key T) {
	s.InMap[key] = 1
}

func (s *Set[T]) SetList(ctx context.Context, keys []T) {
	for _, k := range keys {
		s.InMap[k] = 1
	}
}

func (s *Set[T]) IsSubSet(ctx context.Context, set *Set[T]) bool {
	for v, _ := range set.InMap {
		if _, ok := s.InMap[v]; !ok {
			return false
		}
	}
	return true
}
