package main

import (
	"time"

	"github.com/attapon-th/go-pkg/logger"
	"github.com/robfig/cron/v3"
)

// Rotate file log Every Day
func SetCronJobFileRotaion(loggerWriter *logger.FileWriter) (*cron.Cron, error) {
	runner := cron.New(cron.WithSeconds(), cron.WithLocation(time.Local))
	_, err := runner.AddFunc("0 0 * * * *", func() { loggerWriter.Rotate() })
	return runner, err
}

func main() {
	log.DefaultLogger = logger.
		logger.Debug().Msg("Debug")
	logger.Info().Msg("Info")
	logger.Warn().Msg("Warning")
	logger.Error().Msg("Error")
	logger.Fatal().Msg("Fatal and Exit status 255")
}
