package game

import "time"

// LogType represents the type of log entry.
type LogType string

const (
	LogTypeInfo    = LogType("INFO")
	LogTypeWarning = LogType("WARNING")
	LogTypeDanger  = LogType("DANGER")
	LogTypeSuccess = LogType("SUCCESS")
)

// LogEntry represents a log entry.
type LogEntry struct {
	Time    time.Time
	Type    LogType
	Message string
}
