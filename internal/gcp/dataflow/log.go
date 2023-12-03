package dataflow

import "time"

// A Google log entry
type LogEntry struct {
	Text string
	Time time.Time
}
