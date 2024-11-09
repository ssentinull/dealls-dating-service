package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Trace(v ...interface{})
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})
	Panic(v ...interface{})
}

type Log struct {
	Logger
}

type Options struct{}

func Init() Log {
	return Log{
		Logger: logrus.New(),
	}
}

func (l Log) Trace(v ...interface{}) {
	l.Logger.Trace(v...)
}

func (l Log) Debug(v ...interface{}) {
	l.Logger.Debug(v...)
}

func (l Log) Info(v ...interface{}) {
	l.Logger.Info(v...)
}

func (l Log) Warn(v ...interface{}) {
	l.Logger.Warn(v...)
}

func (l Log) Error(v ...interface{}) {
	l.Logger.Error(v...)
}

func (l Log) Fatal(v ...interface{}) {
	l.Logger.Fatal(v...)
}

func (l Log) Panic(v ...interface{}) {
	l.Logger.Panic(v...)
}
