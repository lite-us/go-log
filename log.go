// Package log is the logging library used by IPFS
// (https://github.com/ipfs/go-ipfs). It uses a modified version of
// https://godoc.org/github.com/whyrusleeping/go-logging .
package log

import (
	"time"

	"go.uber.org/zap"
)

// StandardLogger provides API compatibility with standard printf loggers
// eg. go-logging
type StandardLogger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
}

// EventLogger extends the StandardLogger interface to allow for log items
// containing structured metadata
type EventLogger interface {
	StandardLogger
}

// Logger retrieves an event logger by name
func Logger(system string) *ZapEventLogger {
	if len(system) == 0 {
		setuplog := getLogger("setup-logger")
		setuplog.Error("Missing name parameter")
		system = "undefined"
	}

	return getLogger(system)
}

// ZapEventLogger implements the EventLogger and wraps a zap Sugared Logger.
type ZapEventLogger struct {
	*zap.SugaredLogger
	system string
}

// SetFieldsOnLogger adds the provided key value args as fields to
// the embedded zap.SugaredLogger. These fields are separate to those
// that are provided to the logger via SetFieldsOnAllLoggers and these
// only last for the life time of this particular ZapEventLogger instance.
// Note: the fields will be passed to any children loggers of this logger.
func (zel *ZapEventLogger) SetFieldsOnLogger(args ...interface{}) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	newSugaredLogger := zel.With(args...)
	loggers[zel.system].SugaredLogger = newSugaredLogger
	//zel.SugaredLogger = newSugaredLogger
}

// FormatRFC3339 returns the given time in UTC with RFC3999Nano format.
func FormatRFC3339(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}
