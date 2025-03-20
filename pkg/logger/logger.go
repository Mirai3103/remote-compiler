package logger

import (
	"os"

	prettyconsole "github.com/thessem/zap-prettyconsole"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func init() {
	if os.Getenv("ENV") == "development" {
		log = prettyconsole.NewLogger(zap.DebugLevel)
		return
	}
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.StacktraceKey = ""

	var err error
	log, err = config.Build()
	if err != nil {
		panic(err)
	}
}

func GetLogger() *zap.Logger {
	return log
}
