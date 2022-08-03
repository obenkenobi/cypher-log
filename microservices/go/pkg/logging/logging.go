package logging

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	log "github.com/sirupsen/logrus"
)

// ConfigureGlobalLogging Configure settings based on the app environment on how
// the global logger should work.
func ConfigureGlobalLogging() {
	if environment.IsProduction() {
		log.SetFormatter(getJsonFormatter())
		log.SetLevel(log.InfoLevel)
	} else if environment.IsStaging() {
		log.SetFormatter(getJsonFormatter())
		log.SetLevel(log.DebugLevel)
	} else if environment.IsDevelopment() {
		log.SetFormatter(getTextFormatter())
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetFormatter(getTextFormatter())
		log.SetLevel(log.DebugLevel)
	}
}

// NewLogger Creates a new logger with settings based on the app environment on
// how the logger should work.
func NewLogger() *log.Logger {
	logger := log.New()
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
	return logger
}

func getTextFormatter() *log.TextFormatter {
	return &log.TextFormatter{DisableColors: false, FullTimestamp: true}
}

func getJsonFormatter() *log.JSONFormatter { return &log.JSONFormatter{} }
