package dataflow

import "time"

type MessageLevel string

const (
	LEVEL_UNKNOWN  MessageLevel = "JOB_MESSAGE_IMPORTANCE_UNKNOWN"
	LEVEL_DEBUG    MessageLevel = "JOB_MESSAGE_DEBUG"
	LEVEL_DETAILED MessageLevel = "JOB_MESSAGE_DETAILED"
	LEVEL_BASIC    MessageLevel = "JOB_MESSAGE_BASIC"
	LEVEL_WARNING  MessageLevel = "JOB_MESSAGE_WARNING"
	LEVEL_ERROR    MessageLevel = "JOB_MESSAGE_ERROR"
)

func MessageLevelFromString(in string) MessageLevel {
	switch in {
	case "JOB_MESSAGE_IMPORTANCE_UNKNOWN":
		return LEVEL_UNKNOWN
	case "JOB_MESSAGE_DEBUG":
		return LEVEL_DEBUG
	case "JOB_MESSAGE_DETAILED":
		return LEVEL_DETAILED
	case "JOB_MESSAGE_BASIC":
		return LEVEL_BASIC
	case "JOB_MESSAGE_WARNING":
		return LEVEL_WARNING
	case "JOB_MESSAGE_ERROR":
		return LEVEL_ERROR
	default:
		return LEVEL_UNKNOWN
	}
}

// A Stackdriver log message
type LogMessage struct {
	// The message of the log entry
	Text string

	// The level of the log entry
	Level MessageLevel

	// The time when the log entry was written
	Time time.Time
}
