package src

import "github.com/charmbracelet/log"

func logError(msg string, args ...interface{}) {
	log.Errorf("%s: %v", msg, args)
}

func logFatal(msg string, args ...interface{}) {
	log.Fatalf("%s: %v", msg, args)
}

func logWarn(msg string, args ...interface{}) {
	log.Warn("%s: %v", msg, args)
}

func logInfo(msg string, args ...interface{}) {
	log.Infof("%s: %v", msg, args)
}
