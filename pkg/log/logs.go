package log

import (
	"fmt"
	"log/slog"
	"os"
)

const (
	LevelTrace     = slog.Level(-8)
	LevelDebug     = slog.LevelDebug
	LevelInfo      = slog.LevelInfo
	LevelNotice    = slog.Level(2)
	LevelWarning   = slog.LevelWarn
	LevelError     = slog.LevelError
	LevelEmergency = slog.Level(12)
)

// CustomLogLevel - For details, see https://pkg.go.dev/log/slog#example-HandlerOptions-CustomLevels
func CustomLogLevel(groups []string, a slog.Attr) slog.Attr {
	// Remove time from the output for predictable test output.
	if a.Key == slog.TimeKey {
		return slog.Attr{}
	}

	// Customize the name of the level key and the output string, including
	// custom level values.
	if a.Key == slog.LevelKey {
		// Rename the level key from "level" to "sev".
		a.Key = "sev"

		// Handle custom level values.
		level := a.Value.Any().(slog.Level)

		// This could also look up the name from a map or other structure, but
		// this demonstrates using a switch statement to rename levels. For
		// maximum performance, the string values should be constants, but this
		// example uses the raw strings for readability.
		switch {
		case level < LevelDebug:
			a.Value = slog.StringValue("TRACE")
		case level < LevelInfo:
			a.Value = slog.StringValue("DEBUG")
		case level < LevelNotice:
			a.Value = slog.StringValue("INFO")
		case level < LevelWarning:
			a.Value = slog.StringValue("NOTICE")
		case level < LevelError:
			a.Value = slog.StringValue("WARNING")
		case level < LevelEmergency:
			a.Value = slog.StringValue("ERROR")
		default:
			a.Value = slog.StringValue("EMERGENCY")
		}
	}

	return a
}

type Logger struct {
	l *slog.Logger
}

func (l *Logger) NewJSONLogger(opts *slog.HandlerOptions) *slog.Logger {
	l.l = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	return l.l
}

func New(opts *slog.HandlerOptions) *Logger {
	return &Logger{
		l: slog.New(slog.NewJSONHandler(os.Stdout, opts)),
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
