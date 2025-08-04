package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"calcr/internal/calculator"
	"calcr/internal/handler"
	"calcr/internal/storage"
	"calcr/internal/tracker"

	_ "github.com/go-sql-driver/mysql" // MySQL-Treiber
)

// Konfigurationsvariablen
var (
	useDB        = os.Getenv("USE_DB") == "true"
	mysqlDSN     = os.Getenv("MYSQL_DSN")
	influxURL    = os.Getenv("INFLUX_URL")
	influxToken  = os.Getenv("INFLUX_TOKEN")
	influxOrg    = os.Getenv("INFLUX_ORG")
	influxBucket = os.Getenv("INFLUX_BUCKET")
)

func main() {
	// Zerolog konfigurieren
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("starting webservice...")

	// Initialisierung des Storages
	var resultStorage storage.ResultStorage
	var db *sql.DB
	if useDB {
		log.Info().Msg("storage to a mySQL database enabled.")
		var err error
		db, err = sql.Open("mysql", mysqlDSN)
		if err != nil {
			log.Fatal().Err(err).Msg("error while connecting to database.")
		}
		defer db.Close()

		if err = db.Ping(); err != nil {
			log.Fatal().Err(err).Msg("mysql ping error.")
		}

		resultStorage = storage.NewMySQLStorage(db)
	} else {
		log.Info().Msg("inmemory storage enabled.")
		resultStorage = storage.NewInMemoryStorage()
	}

	// Initialisierung des Trackers
	trackerClient := tracker.NewInfluxDBTracker(influxURL, influxToken, influxOrg, influxBucket)

	r := chi.NewRouter()

	r.Get("/add", handler.HandleAddOperation(calculator.Add, resultStorage, trackerClient))
	r.Get("/subtract", handler.HandleSubOperation(calculator.Subtract, resultStorage, trackerClient))
	r.Get("/multiply", handler.HandleMultiOperation(calculator.Multiply, resultStorage, trackerClient))
	r.Get("/divide", handler.HandleDivOperation(calculator.Divide, resultStorage, trackerClient))
	r.Get("/results", handler.HandleResults(resultStorage))

	port := ":8080"
	log.Info().Msgf("webservice running on http://localhost%s", port)
	http.ListenAndServe(port, r)
}
