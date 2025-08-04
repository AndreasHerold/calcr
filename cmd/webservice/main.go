package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/go-sql-driver/mysql" // MySQL-Treiber
	"github.com/your-username/webservice/internal/calculator"
	"github.com/your-username/webservice/internal/storage"
	"github.com/your-username/webservice/internal/tracker"
)

// Konfigurationsvariablen
var (
	useDB       = os.Getenv("USE_DB") == "true"
	mysqlDSN    = os.Getenv("MYSQL_DSN")
	influxURL   = os.Getenv("INFLUX_URL")
	influxToken = os.Getenv("INFLUX_TOKEN")
	influxOrg   = os.Getenv("INFLUX_ORG")
	influxBucket= os.Getenv("INFLUX_BUCKET")
)

func main() {
	// Zerolog konfigurieren
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("Webservice wird gestartet...")

	// Initialisierung des Storages
	var resultStorage storage.ResultStorage
	var db *sql.DB
	if useDB {
		log.Info().Msg("Speicherung in MySQL-Datenbank aktiviert.")
		var err error
		db, err = sql.Open("mysql", mysqlDSN)
		if err != nil {
			log.Fatal().Err(err).Msg("Verbindung zur MySQL-Datenbank fehlgeschlagen.")
		}
		defer db.Close()

		if err = db.Ping(); err != nil {
			log.Fatal().Err(err).Msg("MySQL-Ping fehlgeschlagen.")
		}

		resultStorage = storage.NewMySQLStorage(db)
	} else {
		log.Info().Msg("Speicherung im Arbeitsspeicher aktiviert.")
		resultStorage = storage.NewInMemoryStorage()
	}

	// Initialisierung des Trackers
	trackerClient := tracker.NewInfluxDBTracker(influxURL, influxToken, influxOrg, influxBucket)

	r := chi.NewRouter()

	r.Get("/add", handleOperation(calculator.Add, "add", resultStorage, trackerClient))
	r.Get("/subtract", handleOperation(calculator.Subtract, "subtract", resultStorage, trackerClient))
	r.Get("/multiply", handleOperation(calculator.Multiply, "multiply", resultStorage, trackerClient))
	r.Get("/divide", handleOperation(calculator.Divide, "divide", resultStorage, trackerClient))
	r.Get("/results", handleResults(resultStorage))

	port := ":8080"
	log.Info().Msgf("Webservice läuft auf http://localhost%s", port)
	http.ListenAndServe(port, r)
}

// handleOperation ist ein Handler-Wrapper für die mathematischen Routen.
func handleOperation(op func(float64, float64) float64, opName string, s storage.ResultStorage, t *tracker.InfluxDBTracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		aStr := r.URL.Query().Get("SummandOne")
		bStr := r.URL.Query().Get("SummandTwo")

		a, errA := strconv.ParseFloat(aStr, 64)
		b, errB := strconv.ParseFloat(bStr, 64)

		if errA != nil || errB != nil {
			log.Error().Err(fmt.Errorf("Fehler beim Parsen der Parameter: a=%s, b=%s", aStr, bStr)).Msg("Ungültige Parameter")
			http.Error(w, "Ungültige Parameter: 'SummandOne' und 'SummandTwo' müssen gültige Zahlen sein.", http.StatusBadRequest)
			return
		}

		start := time.Now()
		result := op(a, b)
		duration := time.Since(start)

		// Ergebnis speichern
		s.AddResult(storage.Result{
			Operation: opName,
			InputA:    a,
			InputB:    b,
			Output:    result,
			Timestamp: time.Now(),
		})

		// Metrik an InfluxDB senden
		go t.TrackOperation(opName, duration, true)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"operation": "%s", "SummandOne": %f, "SummandTwo": %f, "result": %f}`, opName, a, b, result)
	}
}

// handleResults gibt die letzten N Ergebnisse zurück.
func handleResults(s storage.ResultStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recentNStr := r.URL.Query().Get("RecentN")
		recentN, err := strconv.Atoi(recentNStr)
		if err != nil || recentN < 1 || recentN > 20 {
			recentN = 5 // Standardwert
		}

		results, err := s.GetRecentResults(recentN)
		if err != nil {
			log.Error().Err(err).Msg("Fehler beim Abrufen der Ergebnisse")
			http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		for _, res := range results {
			fmt.Fprintf(w, "%s\n", res.String())
		}
	}
}