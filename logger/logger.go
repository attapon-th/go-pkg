package logger

import (
	"os"
	"path/filepath"

	"github.com/phuslu/log"
)

const (
	DEFAULT_LOGFILE = "logs/log.log"
)

var (
	DefaultFileWriter = log.FileWriter{
		FileMode:     0644,
		MaxSize:      1000 * 1024 * 1024, // 1000MB
		EnsureFolder: true,
		LocalTime:    true,
		MaxBackups:   90, // 90 day
		TimeFormat:   "2006-01-02",
	}
	DefaultConsoleJson  = log.IOWriter{os.Stdout}
	DefaultConsoleColor = log.ConsoleWriter{
		ColorOutput:    true,
		EndWithMessage: true,
	}
)

// for set log.Logger default
func SetDefaultlogger(setLogger log.Logger) {
	log.DefaultLogger = setLogger
}

// console log with text format
func GetLogger(logLevel log.Level) log.Logger {
	return log.Logger{
		Level:  log.DebugLevel,
		Caller: 1,
		Writer: &DefaultConsoleColor,
	}
}

// console log with json format
func GetLoggerJson(logLevel log.Level) log.Logger {
	return log.Logger{
		Level:     log.DebugLevel,
		Caller:    1,
		TimeField: "timestamp",
		Writer:    &DefaultConsoleJson,
	}
}

// log file with json format
func GetLoggerFile(filelogName string, logLevel log.Level) log.Logger {
	er := os.MkdirAll(filepath.Dir(filelogName), 0755)
	if er != nil {
		log.Fatal().Err(er)
		os.Exit(1)
	}
	log.Debug().Msgf("logfile: %s", filelogName)
	DefaultFileWriter.Filename = filelogName
	logger := log.Logger{
		Level:  logLevel,
		Writer: &DefaultFileWriter,
	}
	return logger
}

// log file with json format
func GetLoggerFileAndConsole(filelogName string, logLevel log.Level) log.Logger {
	er := os.MkdirAll(filepath.Dir(filelogName), 0755)
	if er != nil {
		log.Fatal().Err(er)
		os.Exit(1)
	}
	log.Debug().Msgf("logfile: %s", filelogName)
	DefaultFileWriter.Filename = filelogName
	logger := log.Logger{
		Level: logLevel,
		Writer: &log.MultiWriter{
			InfoWriter:    &DefaultFileWriter,
			ConsoleWriter: &DefaultConsoleJson,
			ConsoleLevel:  logLevel,
		},
	}
	return logger
}
