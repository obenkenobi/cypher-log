package logger

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	log "github.com/sirupsen/logrus"
)

func getTextFormatter() *log.TextFormatter {
	return &log.TextFormatter{DisableColors: false, FullTimestamp: true}
}

func getJsonFormatter() *log.JSONFormatter { return &log.JSONFormatter{} }

func configureLoggerFromEnv(logger *log.Logger) {
	if environment.IsProduction() {
		logger.SetFormatter(getJsonFormatter())
		logger.SetLevel(log.InfoLevel)
	} else if environment.IsStaging() {
		logger.SetFormatter(getJsonFormatter())
		logger.SetLevel(log.DebugLevel)
	} else if environment.IsDevelopment() {
		logger.SetFormatter(getTextFormatter())
		logger.SetLevel(log.DebugLevel)
	} else {
		logger.SetFormatter(getTextFormatter())
		logger.SetLevel(log.DebugLevel)
	}
}

var Log = log.New()

func ConfigureLoggerFromEnv() { configureLoggerFromEnv(Log) }
