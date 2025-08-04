package tracker

import "time"

type NullTracker struct{}

func NewNullTracker() *NullTracker {
	return &NullTracker{}
}
func (n *NullTracker) TrackOperation(opName string, duration time.Duration, success bool) {
	// No operation tracking in null tracker
}
