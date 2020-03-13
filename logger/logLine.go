//go:generate protoc -I=proto -I=$GOPATH/src -I=$GOPATH/src/github.com/gogo/protobuf/protobuf  --gogoslick_out=. logLineMessage.proto
package logger

import (
	"time"
)

// LogLine is the structure used to hold a log line
type LogLine struct {
	Message   string
	LogLevel  LogLevel
	Args      []interface{}
	Timestamp time.Time
}

func newLogLine(message string, logLevel LogLevel, args ...interface{}) *LogLine {
	return &LogLine{
		Message:   message,
		LogLevel:  logLevel,
		Args:      args,
		Timestamp: time.Now(),
	}
}

// LogLineWrapper is a wrapper over protobuf.LogLineMessage that enables the structure to be used with
// protobuf marshaller
type LogLineWrapper struct {
	LogLineMessage
}

// IsInterfaceNil returns true if there is no value under the interface
func (llw *LogLineWrapper) IsInterfaceNil() bool {
	return llw == nil
}
