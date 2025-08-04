package tracker

import (
	"context"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/rs/zerolog/log"
)

// InfluxDBTracker sendet Metriken an eine InfluxDB-Instanz.
type InfluxDBTracker struct {
	client influxdb2.Client
	org    string
	bucket string
}

// NewInfluxDBTracker erstellt eine neue InfluxDBTracker-Instanz.
func NewInfluxDBTracker(url, token, org, bucket string) *InfluxDBTracker {
	client := influxdb2.NewClient(url, token)
	return &InfluxDBTracker{
		client: client,
		org:    org,
		bucket: bucket,
	}
}

// TrackOperation sendet die Ausführungsdauer einer Operation an InfluxDB.
func (t *InfluxDBTracker) TrackOperation(opName string, duration time.Duration, success bool) {
	writeAPI := t.client.WriteAPIBlocking(t.org, t.bucket)
	p := influxdb2.NewPoint("operation_duration",
		map[string]string{"operation": opName, "status": "success"},
		map[string]interface{}{"duration_ms": duration.Milliseconds()},
		time.Now())

	err := writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		log.Error().Err(err).Msg("Fehler beim Senden der Metrik an InfluxDB")
	}
}