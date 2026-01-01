package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	WithRequestID(requestID string) Logger
}

type DefaultLogger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
	requestID  string
}

func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		debugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		fatalLogger: log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *DefaultLogger) WithRequestID(requestID string) Logger {
	newLogger := *l
	newLogger.requestID = requestID
	return &newLogger
}

func (l *DefaultLogger) formatMessage(message string) string {
	if l.requestID != "" {
		return fmt.Sprintf("[reqID: %s] %s", l.requestID, message)
	}
	return message
}

func (l *DefaultLogger) Debug(args ...interface{}) {
	l.debugLogger.Println(l.formatMessage(fmt.Sprint(args...)))
}

func (l *DefaultLogger) Info(args ...interface{}) {
	l.infoLogger.Println(l.formatMessage(fmt.Sprint(args...)))
}

func (l *DefaultLogger) Warn(args ...interface{}) {
	l.warnLogger.Println(l.formatMessage(fmt.Sprint(args...)))
}

func (l *DefaultLogger) Error(args ...interface{}) {
	l.errorLogger.Println(l.formatMessage(fmt.Sprint(args...)))
}

func (l *DefaultLogger) Fatal(args ...interface{}) {
	l.fatalLogger.Println(l.formatMessage(fmt.Sprint(args...)))
	os.Exit(1)
}

func (l *DefaultLogger) Debugf(format string, args ...interface{}) {
	l.debugLogger.Println(l.formatMessage(fmt.Sprintf(format, args...)))
}

func (l *DefaultLogger) Infof(format string, args ...interface{}) {
	l.infoLogger.Println(l.formatMessage(fmt.Sprintf(format, args...)))
}

func (l *DefaultLogger) Warnf(format string, args ...interface{}) {
	l.warnLogger.Println(l.formatMessage(fmt.Sprintf(format, args...)))
}

func (l *DefaultLogger) Errorf(format string, args ...interface{}) {
	l.errorLogger.Println(l.formatMessage(fmt.Sprintf(format, args...)))
}

func (l *DefaultLogger) Fatalf(format string, args ...interface{}) {
	l.fatalLogger.Println(l.formatMessage(fmt.Sprintf(format, args...)))
	os.Exit(1)
}

var globalLogger Logger = NewDefaultLogger()

func SetLogger(logger Logger) {
	globalLogger = logger
}

func GetLogger() Logger {
	return globalLogger
}

func Debug(args ...interface{}) {
	globalLogger.Debug(args...)
}

func Info(args ...interface{}) {
	globalLogger.Info(args...)
}

func Warn(args ...interface{}) {
	globalLogger.Warn(args...)
}

func Error(args ...interface{}) {
	globalLogger.Error(args...)
}

func Fatal(args ...interface{}) {
	globalLogger.Fatal(args...)
}

func Debugf(format string, args ...interface{}) {
	globalLogger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	globalLogger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	globalLogger.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	globalLogger.Fatalf(format, args...)
}

func WithRequestID(requestID string) Logger {
	return globalLogger.WithRequestID(requestID)
}