package eslog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

// LevelFatal the constant to represent the fatal log level
const (
	LevelFatal = slog.Level(12)
)

// levelNames maps the LevelFatal to string "FATAL"
var levelNames = map[slog.Leveler]string{
	LevelFatal: "FATAL",
}

// eSlogLogger is used to extend slog.
type eSlogLogger struct {
	*slog.Logger
}

// Logger is the default logger which extends slog.
var Logger *eSlogLogger

// logLevel holds the log level of the logger.
var logLevel = &slog.LevelVar{}

func init() {
	initLogger(os.Stdout, false)
}

// initLogger initializes the Logger and enables LevelFatal. Also it sets the default log
// level to LevelDebug
func initLogger(writer io.Writer, overwrite bool) {
	if Logger == nil || overwrite {
		logLevel.Set(slog.LevelDebug)

		opts := &slog.HandlerOptions{
			Level: logLevel,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.LevelKey {
					level := a.Value.Any().(slog.Level)
					levelLabel, exists := levelNames[level]
					if !exists {
						levelLabel = level.String()
					}
					a.Value = slog.StringValue(levelLabel)
				}
				return a
			},
		}

		Logger = &eSlogLogger{
			Logger: slog.New(slog.NewTextHandler(writer, opts)),
		}
	}
}

// SetOutput can be used to overwrite the default output writer os.Stdout. Can be used for
// testing purposes or to swich logging to os.Stderr. In fact that function reinitializes
// the Logger.
func (l *eSlogLogger) SetOutput(w io.Writer) {
	initLogger(w, true)
}

// SetLogLevel sets the LogLevel of the Logger
func (l *eSlogLogger) SetLogLevel(lvl string) error {
	return logLevel.UnmarshalText([]byte(lvl))
}

// Debugf logs at [LevelDebug]. Multiple args are joined with "  ".
func Debug(args ...any) {
	Logger.Debug(strings.Join(convertAnyToString(args...), " "))
}

// Debugf logs at [LevelDebug]. The function uses fmt.Sprintf with given format and args
// and log it.
func Debugf(format string, args ...any) {
	Logger.Debugf(format, args...)
}

// Debugf logs at [LevelDebug]. The function uses fmt.Sprintf with given format and args
// and log it.
func (l eSlogLogger) Debugf(format string, args ...any) {
	l.Debug(fmt.Sprintf(format, args...))
}

// Infof logs at [LevelInfo]. Multiple args are joined with "  ".
func Info(args ...any) {
	Logger.Info(strings.Join(convertAnyToString(args...), " "))
}

// Infof logs at [LevelInfo]. The function uses fmt.Sprintf with given format and args
// and log it.
func Infof(format string, args ...any) {
	Logger.Infof(format, args...)
}

// Infof logs at [LevelInfo]. The function uses fmt.Sprintf with given format and args
// and log it.
func (l eSlogLogger) Infof(format string, args ...any) {
	l.Info(fmt.Sprintf(format, args...))
}

// Warn logs at [LevelWarn]. Multiple args are joined with "  ".
func Warn(args ...any) {
	Logger.Warn(strings.Join(convertAnyToString(args...), " "))
}

// Fatalf logs at [LevelWarn]. The function uses fmt.Sprintf with given format and args
// and log it.
func Warnf(format string, args ...any) {
	Logger.Warnf(format, args...)
}

// Fatalf logs at [LevelWarn]. The function uses fmt.Sprintf with given format and args
// and log it.
func (l eSlogLogger) Warnf(format string, args ...any) {
	l.Warn(format, args...)
}

// Error logs at [LevelError]. Multiple args are joined with "  ".
func Error(args ...any) {
	Logger.Error(strings.Join(convertAnyToString(args...), " "))
}

// Errorf logs at [LevelError]. The function uses fmt.Sprintf with given format and args
// and log it.
func Errorf(format string, args ...any) {
	Logger.Errorf(format, args...)
}

// Errorf logs at [LevelError]. The function uses fmt.Sprintf with given format and args
// and log it.
func (l eSlogLogger) Errorf(format string, args ...any) {
	Logger.Error(fmt.Sprintf(format, args...))
}

// Fatal logs at [LevelFatal].  Multiple args are joined with " ".
func Fatal(args ...any) {
	Logger.Fatal(strings.Join(convertAnyToString(args...), " "))
}

// Fatalf logs at [LevelFatal]. The function uses fmt.Sprintf with given format and args
// and log it. Also it calls os.Exit(1).
func Fatalf(format string, args ...any) {
	Logger.Fatalf(format, args...)
}

// Fatal logs at [LevelFatal]. Also it calls os.Exit(1).
func (l eSlogLogger) Fatal(msg string, args ...any) {
	l.Log(context.Background(), LevelFatal, msg, args...)
	os.Exit(1)
}

// Fatalf logs at [LevelFatal]. The function uses fmt.Sprintf with given format and args
// and log it. Also it calls os.Exit(1).
func (l eSlogLogger) Fatalf(format string, args ...any) {
	l.Fatal(fmt.Sprintf(format, args...))
}

// LogIfError check the given error. If error is nil nothing is logged. If error is not
// nil the loggerFunc is used to log the args. The error is not automatically added to
// args.
func LogIfError(err error, loggerFunc func(args ...any), args ...any) {
	if err != nil {
		loggerFunc(args...)
	}
}

// LogIfErrorf checks the given error. If error is nil nothing is logged. If error is not
// nil the loggerFuncf is used with the given format to print the given args. The error is
// not automatically added to args.
func LogIfErrorf(err error, loggerFuncf func(format string, args ...any), format string, args ...any) {
	if err != nil {
		loggerFuncf(format, args...)
	}
}

// convertAnyToString converts all args of type any to string and
// returns them as []string
func convertAnyToString(args ...any) (strArr []string) {
	strArr = []string{}

	for _, something := range args {
		strArr = append(strArr, fmt.Sprintf("%v", something))
	}

	return strArr
}
