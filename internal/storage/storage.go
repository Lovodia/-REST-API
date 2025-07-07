package storage

import (
	"sync"
)

type ResultStore struct {
	mu      sync.Mutex
	results map[string]map[string]float64
}

func NewResultStore() *ResultStore {
	return &ResultStore{
		results: make(map[string]map[string]float64),
	}
}

func (s *ResultStore) Save(token string, key string, value float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.results[token]; !exists {
		s.results[token] = make(map[string]float64)
	}
	s.results[token][key] = value
}

func (s *ResultStore) GetAllByToken(token string) map[string]float64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	if res, exists := s.results[token]; exists {
		copy := make(map[string]float64, len(res))
		for k, v := range res {
			copy[k] = v
		}
		return copy
	}
	return nil
}
