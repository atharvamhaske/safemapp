// Package safemap provides a small generic map wrapper that is safe for
// concurrent access from multiple goroutines.
//
// It uses sync.RWMutex internally:
//   - reads acquire a shared read lock
//   - writes acquire an exclusive write lock
//
// Use this package when you want a simple thread-safe map API without wiring
// locking logic around every map access in your application code.
//
// Type constraints:
//   - K must be comparable (required by Go maps)
//   - V must be comparable (required for CompareAndSwap)
package safemap

import "sync"

// SafeMap is a generic, concurrency-safe map protected by an RWMutex.
type SafeMap[K comparable, V comparable] struct {
	mu   sync.RWMutex
	data map[K]V
}

// New creates and returns an initialized SafeMap.
func New[K comparable, V comparable]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		data: make(map[K]V),
	}
}

// Set stores v at key k, replacing any existing value.
func (s *SafeMap[K, V]) Set(k K, v V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[k] = v
}

// Get returns the value for key k and whether the key exists.
func (s *SafeMap[K, V]) Get(k K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[k]
	if !ok {
		// Return the type-safe zero value when the key is missing.
		var zero V
		return zero, false
	}
	return val, ok
}

// Delete removes key k from the map.
func (s *SafeMap[K, V]) Delete(k K) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, k)
}

// Len returns the number of entries currently stored.
func (s *SafeMap[K, V]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

// ForEach calls f for each key/value pair.
// Keep callback work lightweight to avoid holding the read lock for long.
func (s *SafeMap[K, V]) ForEach(f func(K, V)) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k, v := range s.data {
		f(k, v)
	}
}

// Exists reports whether key k is present in the map.
func (s *SafeMap[K, V]) Exists(k K) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data[k]
	return ok
}

// CompareAndSwap updates key k from old to new atomically.
// It returns true only when the current value matches old.
func (s *SafeMap[K, V]) CompareAndSwap(k K, old, new V) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, ok := s.data[k]
	if !ok || val != old {
		return false
	}

	s.data[k] = new
	return true
}
