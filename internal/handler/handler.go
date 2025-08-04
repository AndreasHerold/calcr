package handler

import (
	"calcr/internal/storage"
	"calcr/internal/tracker"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	addOperation      = "add"
	subtractOperation = "subtract"
	multiplyOperation = "multiply"
	divideOperation   = "divide"
)

// handleOperation ist ein Handler-Wrapper für die mathematischen Routen.
func HandleAddOperation(op func(float64, float64) float64, s storage.ResultStorage, t tracker.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		aStr := r.URL.Query().Get("SummandOne")
		bStr := r.URL.Query().Get("SummandTwo")

		a, errA := strconv.ParseFloat(aStr, 64)
		b, errB := strconv.ParseFloat(bStr, 64)

		if errA != nil || errB != nil {
			log.Error().Err(fmt.Errorf("error parsing parameters: a=%s, b=%s", aStr, bStr)).Msg("Invalid parameters")
			http.Error(w, "Invalid parameters: 'SummandOne' and 'SummandTwo' must be valid numbers.", http.StatusBadRequest)
			return
		}

		start := time.Now()
		result := op(a, b)
		duration := time.Since(start)

		// Ergebnis speichern
		s.AddResult(storage.Result{
			Operation: addOperation,
			InputA:    a,
			InputB:    b,
			Output:    result,
			Timestamp: time.Now(),
		})

		// Metrik an InfluxDB senden
		go t.TrackOperation(addOperation, duration, true)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"operation": "%s", "SummandOne": %f, "SummandTwo": %f, "Sum": %f}`, addOperation, a, b, result)
	}
}

// handleOperation ist ein Handler-Wrapper für die mathematischen Routen.
func HandleSubOperation(op func(float64, float64) float64, s storage.ResultStorage, t tracker.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		aStr := r.URL.Query().Get("Minuend")
		bStr := r.URL.Query().Get("Subtrahend")

		a, errA := strconv.ParseFloat(aStr, 64)
		b, errB := strconv.ParseFloat(bStr, 64)

		if errA != nil || errB != nil {
			log.Error().Err(fmt.Errorf("error parsing parameters: a=%s, b=%s", aStr, bStr)).Msg("invalid parameters")
			http.Error(w, "Invalid parameters: 'Minuend' and 'Subtrahend' must be valid numbers.", http.StatusBadRequest)
			return
		}

		start := time.Now()
		result := op(a, b)
		duration := time.Since(start)

		// Ergebnis speichern
		s.AddResult(storage.Result{
			Operation: subtractOperation,
			InputA:    a,
			InputB:    b,
			Output:    result,
			Timestamp: time.Now(),
		})

		// Metrik an InfluxDB senden
		go t.TrackOperation(subtractOperation, duration, true)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"operation": "%s", "SummandOne": %f, "SummandTwo": %f, "Difference": %f}`, subtractOperation, a, b, result)
	}
}

// handleOperation ist ein Handler-Wrapper für die mathematischen Routen.
func HandleMultiOperation(op func(float64, float64) float64, s storage.ResultStorage, t tracker.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		aStr := r.URL.Query().Get("FaktorOne")
		bStr := r.URL.Query().Get("FaktorTwo")

		a, errA := strconv.ParseFloat(aStr, 64)
		b, errB := strconv.ParseFloat(bStr, 64)

		if errA != nil || errB != nil {
			log.Error().Err(fmt.Errorf("error parsing parameters: a=%s, b=%s", aStr, bStr)).Msg("invalid parameters")
			http.Error(w, "Invalid parameters: 'FaktorOne' and 'FaktorTwo' must be valid numbers.", http.StatusBadRequest)
			return
		}

		start := time.Now()
		result := op(a, b)
		duration := time.Since(start)

		// Ergebnis speichern
		s.AddResult(storage.Result{
			Operation: multiplyOperation,
			InputA:    a,
			InputB:    b,
			Output:    result,
			Timestamp: time.Now(),
		})

		// Metrik an InfluxDB senden
		go t.TrackOperation(multiplyOperation, duration, true)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"operation": "%s", "SummandOne": %f, "SummandTwo": %f, "Product": %f}`, multiplyOperation, a, b, result)
	}
}

// handleOperation ist ein Handler-Wrapper für die mathematischen Routen.
func HandleDivOperation(op func(float64, float64) float64, s storage.ResultStorage, t tracker.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		aStr := r.URL.Query().Get("Dividend")
		bStr := r.URL.Query().Get("Divisor")

		a, errA := strconv.ParseFloat(aStr, 64)
		b, errB := strconv.ParseFloat(bStr, 64)

		if errA != nil || errB != nil {
			log.Error().Err(fmt.Errorf("error parsing parameters: a=%s, b=%s", aStr, bStr)).Msg("invalid parameters")
			http.Error(w, "Invalid parameters: 'Dividend' and 'Divisor' must be valid numbers.", http.StatusBadRequest)
			return
		}

		start := time.Now()
		result := op(a, b)
		duration := time.Since(start)

		// Ergebnis speichern
		s.AddResult(storage.Result{
			Operation: divideOperation,
			InputA:    a,
			InputB:    b,
			Output:    result,
			Timestamp: time.Now(),
		})

		// Metrik an InfluxDB senden
		go t.TrackOperation(divideOperation, duration, true)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"operation": "%s", "SummandOne": %f, "SummandTwo": %f, "Quotient": %f}`, divideOperation, a, b, result)
	}
}

// handleResults gibt die letzten N Ergebnisse zurück.
func HandleResults(s storage.ResultStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recentNStr := r.URL.Query().Get("RecentN")
		recentN, err := strconv.Atoi(recentNStr)
		if err != nil || recentN < 1 || recentN > 20 {
			recentN = 5 // Standardwert
		}

		results, err := s.GetRecentResults(recentN)
		if err != nil {
			log.Error().Err(err).Msg("error retrieving results")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		for _, res := range results {
			fmt.Fprintf(w, "%s\n", res.String())
		}
	}
}
