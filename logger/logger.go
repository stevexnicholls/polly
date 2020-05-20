package logger

import "errors"

// A global variable so that log functions can be directly accessed
var log Logger

//Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

const (
	//Debug has verbose message
	Debug = "debug"
	//Info is default log level
	Info = "info"
	//Warn is for logging messages about possible issues
	Warn = "warn"
	//Error is for logging errors
	Error = "error"
	//Fatal is for logging fatal messages. The system shutsdown after logging the message.
	Fatal = "fatal"
)

// InstanceZapLogger _
const (
	InstanceZapLogger int = iota
)

var (
	errInvalidLoggerInstance = errors.New("Invalid logger instance")
)

// Logger is our contract for the logger
type Logger interface {
	Debugf(format string, args ...interface{})

	Debugw(msg string, args ...interface{})

	Infof(format string, args ...interface{})

	Infow(msg string, args ...interface{})

	Warnf(format string, args ...interface{})

	Warnw(msg string, args ...interface{})

	Errorf(format string, args ...interface{})

	Errorw(msg string, args ...interface{})

	Fatalf(format string, args ...interface{})

	Fatalw(msg string, args ...interface{})

	Panicf(format string, args ...interface{})

	Panicw(msg string, args ...interface{})

	WithFields(keyValues Fields) Logger
}

// Configuration stores the config for the logger
// For some loggers there can only be one level across writers, for such the level of Console is picked by default
type Configuration struct {
	EnableConsole     bool
	ConsoleJSONFormat bool
	ConsoleLevel      string
	EnableFile        bool
	FileJSONFormat    bool
	FileLevel         string
	FileLocation      string
}

// NewLogger returns an instance of logger
func NewLogger(config Configuration, loggerInstance int) error {
	switch loggerInstance {
	case InstanceZapLogger:
		logger, err := newZapLogger(config)
		if err != nil {
			return err
		}
		log = logger
		return nil

	default:
		return errInvalidLoggerInstance
	}
}

// Debugf _
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Debugw _
func Debugw(msg string, args ...interface{}) {
	log.Debugw(msg, args...)
}

// Infof _
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Infow _
func Infow(msg string, args ...interface{}) {
	log.Infow(msg, args...)
}

// Warnf _
func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Warnw _
func Warnw(msg string, args ...interface{}) {
	log.Warnw(msg, args...)
}

// Errorf _
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Errorw _
func Errorw(msg string, args ...interface{}) {
	log.Errorw(msg, args...)
}

// Fatalf _
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// Fatalw _
func Fatalw(msg string, args ...interface{}) {
	log.Fatalw(msg, args...)
}

// Panicf _
func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

// Panicw _
func Panicw(msg string, args ...interface{}) {
	log.Panicw(msg, args...)
}

// WithFields _
func WithFields(keyValues Fields) Logger {
	return log.WithFields(keyValues)
}
