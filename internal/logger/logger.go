package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type LoggersInterface interface {
	Error(message string, err error)
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Fatal(message string, err error)
	Debug(message string, args ...interface{})
}
type MyLogger struct {
	logger *log.Logger
}

func NewLogger() LoggersInterface {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	myLogger := &MyLogger{
		logger: logger,
	}

	return myLogger
}

func logWithCallerInfo(file string, line int, level string, message string, l *log.Logger, args ...interface{}) {
	caller := fmt.Sprintf("%s:%d", filepath.Base(file), line)
	messageWithCaller := fmt.Sprintf("[%s] %s %s", level, caller, fmt.Sprintf(message, args...))
	l.Println(messageWithCaller)
}

func (l *MyLogger) Error(message string, err error) {
	_, file, line, _ := runtime.Caller(1)
	if l.logger != nil {
		logWithCallerInfo(file, line, "ERROR", "%s: %v", l.logger, message, err)
	} else {
		log.Println("No logger available.")
	}
}

func (l *MyLogger) Info(message string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if l.logger != nil {
		logWithCallerInfo(file, line, "INFO", message, l.logger, args...)
	} else {
		log.Println("No logger available.")
	}
}

func (l *MyLogger) Warn(message string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if l.logger != nil {
		logWithCallerInfo(file, line, "WARN", message, l.logger, args...)
	} else {
		log.Println("No logger available.")
	}
}

func (l *MyLogger) Fatal(message string, err error) {
	_, file, line, _ := runtime.Caller(1)
	if l.logger != nil {
		logWithCallerInfo(file, line, "FATAL", "%s: %v", l.logger, message, err)
		os.Exit(1) // Завершаем приложение с кодом ошибки
	} else {
		log.Println("No logger available.")
	}
}

func (l *MyLogger) Debug(message string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if l.logger != nil {
		logWithCallerInfo(file, line, "DEBUG", message, l.logger, args...)
	} else {
		log.Println("No logger available.")
	}
}
