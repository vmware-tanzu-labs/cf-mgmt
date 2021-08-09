package lo

import (
	"os"

	"github.com/op/go-logging"
)

const (
	LOG_MODULE = "lo.G_logger"
)

var (
	G Logger
)

type Logger interface {
	Critical(args ...interface{})
	Criticalf(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Notice(args ...interface{})
	Noticef(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
}

func init() {

	if logLevel, err := logging.LogLevel(os.Getenv("LOG_LEVEL")); err == nil {
		logging.SetLevel(logLevel, LOG_MODULE)

	} else {
		logging.SetLevel(logging.INFO, LOG_MODULE)
	}
	logging.SetFormatter(logging.GlogFormatter)
	G = logging.MustGetLogger(LOG_MODULE)
}
