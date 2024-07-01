package utils

import (
	"fmt"
	"sync"
)

type SyncMap[K comparable, V interface{}] struct {
	smap map[K]V
	mu   sync.RWMutex
}

func NewSyncMap[K comparable, V interface{}]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		smap: make(map[K]V),
		mu:   sync.RWMutex{},
	}
}

func (m *SyncMap[K, V]) Set(key K, val V) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Exists(key) {
		return fmt.Errorf("data already exists")
	}

	m.smap[key] = val

	return nil
}

func (m *SyncMap[K, V]) Get(key K) (V, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var zero V

	if !m.Exists(key) {
		return zero, ErrNotFound
	}

	return m.smap[key], nil
}

func (m *SyncMap[K, V]) GetAll() ([]V, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var vs []V

	for _, v := range m.smap {
		vs = append(vs, v)
	}

	return vs, nil
}

func (m *SyncMap[K, V]) Remove(key K) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.Exists(key) {
		return fmt.Errorf("data does'nt exist")
	}

	delete(m.smap, key)

	return nil
}

func (m *SyncMap[K, V]) Update(key K, newval V) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.Exists(key) {
		return fmt.Errorf("data does'nt exists")
	}

	m.smap[key] = newval

	return nil
}

func (m *SyncMap[K, V]) Count() int {
	return len(m.smap)
}

func (m *SyncMap[K, V]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k := range m.smap {
		delete(m.smap, k)
	}
}

func (m *SyncMap[K, V]) Exists(key K) bool {
	_, ok := m.smap[key]
	return ok
}
