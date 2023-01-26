package logger

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	log "github.com/sirupsen/logrus"
)

func getTextFormatter() *log.TextFormatter {
	return &log.TextFormatter{
		DisableColors:    false,
		ForceColors:      true,
		FullTimestamp:    true,
		PadLevelText:     true,
		QuoteEmptyFields: true,
		ForceQuote:       true,
	}
}

func getJsonFormatter() *log.JSONFormatter { return &log.JSONFormatter{PrettyPrint: true} }

func configureLoggerFromEnv(logger *log.Logger) {
	logger.SetReportCaller(true)
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

func newLogger() *log.Logger {
	logger := log.New()
	logger.SetReportCaller(true)
	logger.SetFormatter(getTextFormatter())
	logger.SetLevel(log.InfoLevel)
	return logger
}

var Log = newLogger()

func ConfigureLoggerFromEnv() { configureLoggerFromEnv(Log) }
