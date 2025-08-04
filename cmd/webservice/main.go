package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"calcr/internal/calculator"
	"calcr/internal/config"
	"calcr/internal/handler"
	"calcr/internal/storage"
	"calcr/internal/tracker"

	_ "github.com/go-sql-driver/mysql" // MySQL-Treiber
)

// Standard-Konfigurationsdatei
const defaultConfigFile = "config/config.yaml"

func main() {
	// Zerolog konfigurieren
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("starting webservice...")

	// Konfiguration laden
	cfg, err := config.LoadConfig(defaultConfigFile)
	if err != nil {
		log.Fatal().Err(err).Msg("error loading configuration")
	}

	// Initialisierung des Storages
	resultStorage := initResultStorage(cfg)

	// Initialisierung des Trackers
	trackerClient := initTracker(cfg)

	r := chi.NewRouter()

	r.Get("/add", handler.HandleAddOperation(calculator.Add, resultStorage, trackerClient))
	r.Get("/subtract", handler.HandleSubOperation(calculator.Subtract, resultStorage, trackerClient))
	r.Get("/multiply", handler.HandleMultiOperation(calculator.Multiply, resultStorage, trackerClient))
	r.Get("/divide", handler.HandleDivOperation(calculator.Divide, resultStorage, trackerClient))
	r.Get("/results", handler.HandleResults(resultStorage))

	port := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Info().Msgf("webservice running on http://localhost%s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal().Err(err).Msg("server error")
	}
}

func initResultStorage(cfg *config.Config) storage.ResultStorage {
	var resultStorage storage.ResultStorage
	if cfg.Database.Enabled {
		log.Info().Msg("storage to a mySQL database enabled.")
		db, err := sql.Open("mysql", cfg.Database.DSN)
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
	return resultStorage
}

func initTracker(cfg *config.Config) tracker.Tracker {
	var trackerClient tracker.Tracker
	if cfg.InfluxDB.Enabled {
		log.Info().Msg("tracker enabled.")
		trackerClient = tracker.NewInfluxDBTracker(
			cfg.InfluxDB.URL,
			cfg.InfluxDB.Token,
			cfg.InfluxDB.Org,
			cfg.InfluxDB.Bucket)
	} else {
		log.Info().Msg("tracker disabled, using null tracker.")
		trackerClient = tracker.NewNullTracker()
	}
	return trackerClient
}
