package extendedslog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"reflect"
	"slices"
)

const (
	LevelFatal = slog.Level(12)
)

var LevelNames = map[slog.Leveler]string{
	LevelFatal: "FATAL",
}

type ExtendedSlogLogger struct {
	*slog.Logger
}

var Logger *ExtendedSlogLogger

var LogLevel = &slog.LevelVar{}

func init() {
	InitLogger(os.Stdout, false)
}

func InitLogger(writer io.Writer, overwrite bool) {
	if Logger == nil || overwrite {
		LogLevel.Set(slog.LevelDebug)

		opts := &slog.HandlerOptions{
			Level: LogLevel,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.LevelKey {
					level := a.Value.Any().(slog.Level)
					levelLabel, exists := LevelNames[level]
					if !exists {
						levelLabel = level.String()
					}
					a.Value = slog.StringValue(levelLabel)
				}
				return a
			},
		}

		Logger = &ExtendedSlogLogger{
			slog.New(slog.NewTextHandler(writer, opts)),
		}
	}
}

func (l *ExtendedSlogLogger) SetOutput(w io.Writer) {
	InitLogger(w, true)
}

func (l *ExtendedSlogLogger) SetLogLevel(lvl string) error {
	return LogLevel.UnmarshalText([]byte(lvl))
}

func (l ExtendedSlogLogger) Debugf(format string, args ...any) {
	l.Debug(fmt.Sprintf(format, args...))
}

func (l ExtendedSlogLogger) Infof(format string, args ...any) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l ExtendedSlogLogger) Warnf(format string, args ...any) {
	l.Warn(format, args...)
}

func (l ExtendedSlogLogger) Error(err error, args ...any) {
	if err != nil {
		l.Logger.Error(err.Error(), args...)
	}
}

func (l ExtendedSlogLogger) Errorf(format string, err error, args ...any) {

	if err != nil {
		args = slices.Insert(args, 0, reflect.ValueOf(err).Interface())
		l.Logger.Error(fmt.Sprintf(format, args...))
	}
}

func (l ExtendedSlogLogger) Fatal(msg string) {
	l.Log(context.Background(), LevelFatal, msg)
	os.Exit(1)
}

func (l ExtendedSlogLogger) Fatalf(format string, err error, args ...any) {
	if err != nil {
		args = slices.Insert(args, 0, reflect.ValueOf(err).Interface())
		l.Fatal(fmt.Sprintf(format, args...))
	}
}
