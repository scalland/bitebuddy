package log

import (
	"fmt"
	"log/slog"
	"os"
)

type Logger struct {
	l *slog.Logger
}

func (l *Logger) NewJSONLogger() *slog.Logger {
	l.l = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return l.l
}

func New() *Logger {
	return &Logger{
		l: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.l.Debug(fmt.Sprintf(format, v...))
}

func (l *Logger) Info(v ...interface{}) {
	l.l.Info(fmt.Sprint(v...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.l.Info(fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(v ...interface{}) {
	l.l.Warn(fmt.Sprint(v...))
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.l.Warn(fmt.Sprintf(format, v...))
}

func (l *Logger) Error(v ...interface{}) {
	l.l.Error(fmt.Sprint(v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.l.Error(fmt.Sprintf(format, v...))
}

func (l *Logger) Fatal(v ...interface{}) {
	l.l.Error(fmt.Sprint(v...))
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.l.Error(fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (l *Logger) Panic(v ...interface{}) {
	l.l.Error(fmt.Sprint(v...))
}

//func (l *Logger) Panicf(format string, v ...interface{})   {}
//func (l *Logger) Debugln(v ...interface{})                 {}
//func (l *Logger) Infoln(v ...interface{})                  {}
//func (l *Logger) Infolnf(format string, v ...interface{})  {}
//func (l *Logger) Warnln(v ...interface{})                  {}
//func (l *Logger) Errorln(v ...interface{})                 {}
//func (l *Logger) Errorlnf(format string, v ...interface{}) {}
//func (l *Logger) Fatalln(v ...interface{})                 {}
//func (l *Logger) Fatallnf(format string, v ...interface{}) {}
//func (l *Logger) Panicln(v ...interface{})                 {}
//func (l *Logger) Paniclnf(format string, v ...interface{}) {}
//func (l *Logger) Debugc(format string, v ...interface{})   {}
//func (l *Logger) Infc(format string, v ...interface{})     {}
//func (l *Logger) Inficc(format string, v ...interface{})   {}
//func (l *Logger) Warnc(format string, v ...interface{})    {}
//func (l *Logger) Errorc(format string, v ...interface{})   {}
//func (l *Logger) Erroric(format string, v ...interface{})  {}
//func (l *Logger) Fatallnc(format string, v ...interface{}) {}
