package reposted

import (
	"iter"
	"sync"
)

type SafeMap[K comparable, V any] struct {
	m map[K]V
	l *sync.RWMutex
}

func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		m: map[K]V{},
		l: &sync.RWMutex{},
	}
}

func (s *SafeMap[K, V]) Get(key K) V {
	s.l.RLock()
	defer s.l.RUnlock()
	return s.m[key]
}

func (s *SafeMap[K, V]) Get2(key K) (V, bool) {
	s.l.RLock()
	defer s.l.RUnlock()
	v, ok := s.m[key]
	return v, ok
}

func (s *SafeMap[K, V]) Set(key K, value V) {
	s.l.Lock()
	s.m[key] = value
	s.l.Unlock()
}

func (s *SafeMap[K, V]) Len() int {
	s.l.RLock()
	defer s.l.RUnlock()
	return len(s.m)
}


func (s *SafeMap[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		s.l.RLock()
		for k, v := range s.m {
			s.l.RUnlock()
			if !yield(k, v) {
				return
			}
			s.l.RLock()
		}
	}
}
