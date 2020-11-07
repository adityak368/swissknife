package logger

import (
	"io"
	"os"

	"github.com/adityak368/rlog"
)

// LogLevel defines the loglevel
type LogLevel string

const (
	// LogLevelCritical is the critical log level
	LogLevelCritical = "CRITICAL"
	// LogLevelError is the error log level
	LogLevelError = "ERROR"
	// LogLevelWarn is the warn log level
	LogLevelWarn = "WARN"
	// LogLevelDebug is the debug log level
	LogLevelDebug = "DEBUG"
	// LogLevelInfo is the info log level
	LogLevelInfo = "INFO"
	// LogLevelTrace is the trace log level
	LogLevelTrace = "TRACE"
)

// tracelevel sets the tracelevel
var tracelevel = 1

// Critical prints a message if log level is set to CRITICAL or lower.
func Critical(a ...interface{}) {
	rlog.Critical(a...)
}

//Criticalf prints a message if log level is set to CRITICAL or lower, with formatting.
func Criticalf(format string, a ...interface{}) {
	rlog.Criticalf(format, a...)
}

//Debug prints a message if log level is set to DEBUG or lower.
func Debug(a ...interface{}) {
	rlog.Debug(a...)
}

//Debugf prints a message if log level is set to DEBUG or lower, with formatting.
func Debugf(format string, a ...interface{}) {
	rlog.Debugf(format, a...)
}

//Error prints a message if log level is set to ERROR or lower.
func Error(a ...interface{}) {
	rlog.Error(a...)
}

//Errorf prints a message if log level is set to ERROR or lower, with formatting.
func Errorf(format string, a ...interface{}) {
	rlog.Errorf(format, a...)
}

//Warn prints a message if log level is set to WARN or lower, with formatting.
func Warn(a ...interface{}) {
	rlog.Warn(a...)
}

//Warnf prints a message if log level is set to WARN or lower, with formatting.
func Warnf(format string, a ...interface{}) {
	rlog.Warnf(format, a...)
}

//Info prints a message if log level is set to INFO or lower.
func Info(a ...interface{}) {
	rlog.Info(a...)
}

//Infof prints a message if log level is set to INFO or lower, with formatting.
func Infof(format string, a ...interface{}) {
	rlog.Infof(format, a...)
}

//Trace prints a message if log level is set to TRACE or lower.
func Trace(a ...interface{}) {
	rlog.Trace(tracelevel, a...)
}

//Tracef prints a message if log level is set to TRACE or lower, with formatting.
func Tracef(format string, a ...interface{}) {
	rlog.Tracef(tracelevel, format, a...)
}

//SetOutput prints a message if log level is set to CRITICAL or lower, with formatting.
func SetOutput(writer io.Writer) {
	rlog.SetOutput(writer)
}

// SetLogLevel sets the loglevel of the logger
func SetLogLevel(level LogLevel) {
	os.Setenv("RLOG_LOG_LEVEL", "DEBUG")
	rlog.UpdateEnv()
}

//SetLogOutputFile sets the logger's output file
func SetLogOutputFile(logFileName string) {
	os.Setenv("RLOG_LOG_FILE", logFileName)
	rlog.UpdateEnv()
}

//SetShowCallerInfo sets if the logger should log the caller info
func SetShowCallerInfo(isEnable bool) {
	os.Setenv("RLOG_CALLER_INFO", boolToStr(isEnable))
	os.Setenv("RLOG_GOROUTINE_ID", boolToStr(isEnable))
	rlog.UpdateEnv()
}

func boolToStr(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
