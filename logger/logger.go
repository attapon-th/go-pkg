package logger

import (
	"os"
	"path/filepath"

	"github.com/phuslu/log"
)

type (
	Logger        = log.Logger
	FileWriter    = log.FileWriter
	MultiWriter   = log.MultiWriter
	ConsoleWriter = log.ConsoleWriter
	Entry         = log.Entry
)

const (
	DEFAULT_LOGFILE = "logs/log.log"
)

var (
	DefaultFileWriter = FileWriter{
		FileMode:     0644,
		MaxSize:      1000 * 1024 * 1024, // 1000MB
		EnsureFolder: true,
		LocalTime:    true,
		MaxBackups:   90, // 90 day
		TimeFormat:   "2006-01-02",
	}
	DefaultConsoleJson  = log.IOWriter{os.Stdout}
	DefaultConsoleColor = ConsoleWriter{
		ColorOutput:    true,
		EndWithMessage: true,
	}
)

func init() {
	SetDefaultlogger(GetLogger(DebugLevel))
}

// for set logger default
func SetDefaultlogger(setLogger Logger) {
	log.DefaultLogger = setLogger
}

// console log with text format
func GetLogger(logLevel Level) Logger {
	return Logger{
		Level:  DebugLevel,
		Caller: 1,
		Writer: &DefaultConsoleColor,
	}
}

// console log with json format
func GetLoggerJson(logLevel Level) Logger {
	return log.Logger{
		Level:     DebugLevel,
		Caller:    1,
		TimeField: "timestamp",
		Writer:    &DefaultConsoleJson,
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
	DefaultFileWriter.Filename = filelogName
	logger := Logger{
		Level:  logLevel,
		Writer: &DefaultFileWriter,
	}
	return logger
}

// log file with json format
func GetLoggerFileAndConsole(filelogName string, logLevel Level) Logger {
	er := os.MkdirAll(filepath.Dir(filelogName), 0755)
	if er != nil {
		log.Fatal().Err(er)
		os.Exit(1)
	}
	Debug().Msgf("logfile: %s", filelogName)
	DefaultFileWriter.Filename = filelogName
	logger := Logger{
		Level: logLevel,
		Writer: &MultiWriter{
			InfoWriter:    &DefaultFileWriter,
			ConsoleWriter: &DefaultConsoleJson,
			ConsoleLevel:  logLevel,
		},
	}
	return logger
}
