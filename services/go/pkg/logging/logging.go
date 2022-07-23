package logging

import log "github.com/sirupsen/logrus"

func GetTextFormatter() *log.TextFormatter {
	return &log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	}
}

func ConfigTextLogging() {
	log.SetFormatter(GetTextFormatter())
	log.SetReportCaller(true)
}

func ConfigTextLoggingWithLogger(logger *log.Logger) {
	logger.SetFormatter(GetTextFormatter())
	logger.SetReportCaller(true)
}
