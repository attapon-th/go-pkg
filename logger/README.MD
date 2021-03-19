# Logger

#### ***Base On: [github.com/phuslu/log](https://github.com/phuslu/log)***

### Install 
```shell
go get -u github.com/attapon-th/go-pkg/logger
```
**Or**
```go
import "github.com/attapon-th/go-pkg/logger"
```

### Signature
```go
import "github.com/attapon-th/go-pkg/logger"

var log logger.Logger

// console log with text format
log = GetLogger(logger.DebugLevel)
log.Info().Msg("console log with text format")

// console log with json format
//func GetLoggerJson(logLevel Level) Logger {

// log file with json format
// func GetLoggerFile(filelogName string, logLevel Level) Logger {
```



### Example

```go
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
	logger.Debug().Msg("Debug")
	logger.Info().Msg("Info")
	logger.Warn().Msg("Warning")
	logger.Error().Msg("Error")
	logger.Fatal().Msg("Fatal and Exit status 255")
}

```


