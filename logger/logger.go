package logger

import (
	"os"
	"path/filepath"

	"github.com/phuslu/log"
)

type (
	Logger        = log.Logger
	IOWriter      = log.IOWriter
	FileWriter    = log.FileWriter
	MultiWriter   = log.MultiWriter
	AsyncWriter   = log.AsyncWriter
	ConsoleWriter = log.ConsoleWriter
	Entry         = log.Entry
)

const (
	DEFAULT_LOGFILE = "logs/log.log"
)

func init() {
	log.DefaultLogger = GetLogger(DebugLevel)
}

// for set logger default
func GetDefaultlogger() *Logger {
	return &log.DefaultLogger
}

// console log with text format
func GetLogger(logLevel Level) Logger {
	return Logger{
		Level:  DebugLevel,
		Caller: 1,
		Writer: &ConsoleWriter{
			ColorOutput:    true,
			EndWithMessage: true,
		},
	}
}

// console log with json format
func GetLoggerJson(logLevel Level) Logger {
	return log.Logger{
		Level:     DebugLevel,
		Caller:    1,
		TimeField: "timestamp",
		Writer:    &IOWriter{os.Stdout},
	}
}

// log file with json format
func GetLoggerFile(filelogName string, logLevel Level) Logger {
	er := os.MkdirAll(filepath.Dir(filelogName), 0755)
	if er != nil {
		log.Fatal().Err(er)
		os.Exit(1)
	}
	log.Debug().Msgf("logfile: %s", filelogName)
	logger := Logger{
		Level: logLevel,
		Writer: &FileWriter{
			Filename:     filelogName,
			FileMode:     0644,
			MaxSize:      500 * 1024 * 1024,
			EnsureFolder: true,
			LocalTime:    true,
		},
	}
	return logger
}
