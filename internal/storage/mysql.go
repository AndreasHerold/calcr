package storage

import (
	"database/sql"
	"time"
)

// MySQLStorage speichert Ergebnisse in einer MySQL-Datenbank.
type MySQLStorage struct {
	db *sql.DB
}

// NewMySQLStorage erstellt eine neue MySQLStorage-Instanz.
func NewMySQLStorage(db *sql.DB) *MySQLStorage {
	return &MySQLStorage{db: db}
}

func (s *MySQLStorage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// AddResult fügt ein Ergebnis in die MySQL-Datenbank ein.
func (s *MySQLStorage) AddResult(result Result) error {
	_, err := s.db.Exec("INSERT INTO results (operation, input_a, input_b, output, timestamp) VALUES (?, ?, ?, ?, ?)",
		result.Operation, result.InputA, result.InputB, result.Output, result.Timestamp)
	return err
}

// GetRecentResults gibt die letzten N Ergebnisse aus der MySQL-Datenbank zurück.
func (s *MySQLStorage) GetRecentResults(n int) ([]Result, error) {
	rows, err := s.db.Query("SELECT operation, input_a, input_b, output, timestamp FROM results ORDER BY timestamp DESC LIMIT ?", n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Result
	var dbTimestamp string
	for rows.Next() {
		var r Result
		if err := rows.Scan(&r.Operation, &r.InputA, &r.InputB, &r.Output, &dbTimestamp); err != nil {
			return nil, err
		}
		r.Timestamp, err = time.Parse("2006-01-02 15:04:05", dbTimestamp)
		if err != nil {
			continue // Fehler beim Parsen des Timestamps, überspringe diesen Eintrag
		}
		results = append(results, r)
	}
	return results, nil
}
