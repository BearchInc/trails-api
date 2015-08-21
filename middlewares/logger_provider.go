package middlewares

import (
	"appengine"
	"github.com/go-martini/martini"
)

// Logger wraps appengine context
// providing only the logging functionality
type Logger struct {
	context appengine.Context
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.context.Infof(format, args...)
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.context.Debugf(format, args...)
}

func (logger *Logger) Warningf(format string, args ...interface{}) {
	logger.context.Warningf(format, args...)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.context.Errorf(format, args...)
}

// LoggerProvider is a martini middleware responsible for providing
// an injectable instance of Logger
func LoggerProvider(c martini.Context, context appengine.Context) {
	c.Map(&Logger{context})
}
