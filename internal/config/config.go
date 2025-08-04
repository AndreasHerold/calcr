package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		Enabled bool   `yaml:"enabled"`
		DSN     string `yaml:"dsn"`
	} `yaml:"database"`

	InfluxDB struct {
		Enabled bool   `yaml:"enabled"`
		URL     string `yaml:"url"`
		Token   string `yaml:"token"`
		Org     string `yaml:"org"`
		Bucket  string `yaml:"bucket"`
	} `yaml:"influxdb"`

	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
}

// LoadConfig lädt die Konfiguration aus der angegebenen YAML-Datei
func LoadConfig(filename string) (*Config, error) {
	config := &Config{}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	// Überschreibe Werte mit Umgebungsvariablen, falls vorhanden
	if os.Getenv("USE_DB") != "" {
		config.Database.Enabled = os.Getenv("USE_DB") == "true"
	}
	if os.Getenv("MYSQL_DSN") != "" {
		config.Database.DSN = os.Getenv("MYSQL_DSN")
	}
	if os.Getenv("INFLUX_URL") != "" {
		config.InfluxDB.URL = os.Getenv("INFLUX_URL")
	}
	if os.Getenv("INFLUX_TOKEN") != "" {
		config.InfluxDB.Token = os.Getenv("INFLUX_TOKEN")
	}
	if os.Getenv("INFLUX_ORG") != "" {
		config.InfluxDB.Org = os.Getenv("INFLUX_ORG")
	}
	if os.Getenv("INFLUX_BUCKET") != "" {
		config.InfluxDB.Bucket = os.Getenv("INFLUX_BUCKET")
	}

	return config, nil
}
