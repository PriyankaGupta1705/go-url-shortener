package store

import (
	"errors"
	"sync"
)

type Store struct {
	mu   sync.RWMutex
	urls map[string]string
}

func NewStore() *Store {
	return &Store{
		urls: make(map[string]string),
	}
}

func (s *Store) Save(code, url string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.urls[code] = url
}

func (s *Store) Get(code string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.urls[code]
	if !ok {
		return "", errors.New("url not found")
	}
	return url, nil
}
