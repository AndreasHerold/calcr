package storage

import (
	"fmt"
	"time"
)

// Result repräsentiert eine gespeicherte Berechnung.
type Result struct {
	ID        int64
	Operation string
	InputA    float64
	InputB    float64
	Output    float64
	Timestamp time.Time
}

// String gibt eine formatierte String-Repräsentation des Ergebnisses zurück.
func (r Result) String() string {
	return fmt.Sprintf("%f %s %f = %f", r.InputA, r.Operation, r.InputB, r.Output)
}

// ResultStorage ist ein Interface für die Speicherung von Ergebnissen.
type ResultStorage interface {
	AddResult(result Result) error
	GetRecentResults(n int) ([]Result, error)
	Close() error
}
