package store

import (
	"fmt"
	"sync"
)

type Storage[T any] struct {
	store map[string]T
	sync.RWMutex
}

func NewStorage[T any]() *Storage[T] {
	return &Storage[T]{
		store: make(map[string]T),
	}
}

func (s *Storage[T]) Put(key string, value T) {
	s.Lock()
	s.store[key] = value
	s.Unlock()
}

func (s *Storage[T]) Get(key string) (*T, error) {
	s.RLock()
	value, ok := s.store[key]
	s.RUnlock()
	if !ok {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	return &value, nil
}

func (s *Storage[T]) Delete(key string) {
	s.Lock()
	delete(s.store, key)
	s.Unlock()
}
