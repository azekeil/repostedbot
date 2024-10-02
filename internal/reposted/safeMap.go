package reposted

import (
	"iter"
	"sync"
)

type SafeMap[K comparable, V any] struct {
	M map[K]V
	l *sync.RWMutex
}

func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		M: map[K]V{},
		l: &sync.RWMutex{},
	}
}

func (s *SafeMap[K, V]) Get(key K) V {
	// if s.l == nil {
	// 	s.l = &sync.RWMutex{}
	// }
	s.l.RLock()
	defer s.l.RUnlock()
	return s.M[key]
}

func (s *SafeMap[K, V]) Get2(key K) (V, bool) {
	// if s.l == nil {
	// 	s.l = &sync.RWMutex{}
	// }
	s.l.RLock()
	defer s.l.RUnlock()
	v, ok := s.M[key]
	return v, ok
}

func (s *SafeMap[K, V]) Set(key K, value V) {
	// if s.l == nil {
	// 	s.l = &sync.RWMutex{}
	// }
	s.l.Lock()
	s.M[key] = value
	s.l.Unlock()
}

func (s *SafeMap[K, V]) Len() int {
	// if s.l == nil {
	// 	s.l = &sync.RWMutex{}
	// }
	s.l.RLock()
	defer s.l.RUnlock()
	return len(s.M)
}

func (s *SafeMap[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		// if s.l == nil {
		// 	s.l = &sync.RWMutex{}
		// }
		s.l.RLock()
		for k, v := range s.M {
			s.l.RUnlock()
			if !yield(k, v) {
				return
			}
			s.l.RLock()
		}
		s.l.RUnlock()
	}
}

func (s *SafeMap[K, V]) EnsureMutex() {
	if s.l == nil {
		s.l = &sync.RWMutex{}
	}
}
