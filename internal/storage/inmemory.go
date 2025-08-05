package storage

import (
	"sync"
)

// InMemoryStorage speichert Ergebnisse im Arbeitsspeicher.
type InMemoryStorage struct {
	mu      sync.Mutex
	results []Result
}

// NewInMemoryStorage erstellt eine neue InMemoryStorage-Instanz.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		results: make([]Result, 0, 20),
	}
}

// Close schließt die InMemoryStorage-Instanz. Da es sich um eine In-Memory-Lösung handelt, ist hier keine Aktion erforderlich.
func (s *InMemoryStorage) Close() error {
	// Keine Aktion erforderlich, da es sich um eine In-Memory-Lösung handelt.
	return nil
}

// AddResult fügt ein Ergebnis hinzu.
func (s *InMemoryStorage) AddResult(result Result) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.results = append(s.results, result)
	return nil
}

// GetRecentResults gibt die letzten N Ergebnisse zurück.
func (s *InMemoryStorage) GetRecentResults(n int) ([]Result, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if n > len(s.results) {
		n = len(s.results)
	}

	return s.results[len(s.results)-n:], nil
}
