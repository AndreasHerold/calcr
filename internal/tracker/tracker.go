package tracker

import "time"

type Tracker interface {
	TrackOperation(opName string, duration time.Duration, success bool)
}
