package logger

import "github.com/phuslu/log"

type Level = log.Level

const (
	// TraceLevel defines trace log level.
	TraceLevel Level = 1
	// DebugLevel defines debug log level.
	DebugLevel Level = 2
	// InfoLevel defines info log level.
	InfoLevel Level = 3
	// WarnLevel defines warn log level.
	WarnLevel Level = 4
	// ErrorLevel defines error log level.
	ErrorLevel Level = 5
	// FatalLevel defines fatal log level.
	FatalLevel Level = 6
	// PanicLevel defines panic log level.
	PanicLevel Level = 7
	// NoLevel defines an absent log level.
	noLevel Level = 8
)
