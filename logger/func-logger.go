package logger

import "github.com/phuslu/log"

func Trace() *log.Entry { return log.Trace() }
func Debug() *log.Entry { return log.Debug() }
func Info() *log.Entry  { return log.Info() }
func Warn() *log.Entry  { return log.Warn() }
func Error() *log.Entry { return log.Error() }
func Fatal() *log.Entry { return log.Fatal() }
