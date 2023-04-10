package game

import "time"

type LogType string

const (
	LogTypeInfo    = LogType("INFO")
	LogTypeWarning = LogType("WARNING")
	LogTypeDanger  = LogType("DANGER")
	LogTypeSuccess = LogType("SUCCESS")
)

type LogEntry struct {
	Time    time.Time
	Type    LogType
	Message string
}
