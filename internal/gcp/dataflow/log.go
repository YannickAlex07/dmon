package dataflow

import "time"

// A Google log entry
type DataflowLogMessage struct {
	Text  string
	Level string
	Time  time.Time
}
