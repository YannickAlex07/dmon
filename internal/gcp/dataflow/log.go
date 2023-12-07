package dataflow

import "time"

// A Google log entry
type DataflowLogEntry struct {
	Text  string
	Level string
	Time  time.Time
}
