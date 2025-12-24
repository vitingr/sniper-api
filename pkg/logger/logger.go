package logger

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	With(keysAndValues ...interface{}) Logger
}

type zerologLogger struct {
	logger zerolog.Logger
}

func parseLogLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

func New(level string) Logger {
	logLevel := parseLogLevel(level)

	var output io.Writer = os.Stdout
	if os.Getenv("ENVIRONMENT") == "production" {
		zerolog.TimeFieldFormat = time.RFC3339
		output = os.Stdout
	} else {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
			NoColor:    false,
		}
	}

	logger := zerolog.New(output).Level(logLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	return &zerologLogger{logger: logger}
}

func (l *zerologLogger) Info(msg string, keysAndValues ...interface{}) {
	event := l.logger.Info()
	l.addFields(event, keysAndValues...)
	event.Msg(msg)
}

func (l *zerologLogger) Error(msg string, keysAndValues ...interface{}) {
	event := l.logger.Error()
	l.addFields(event, keysAndValues...)
	event.Msg(msg)
}

func (l *zerologLogger) Debug(msg string, keysAndValues ...interface{}) {
	event := l.logger.Debug()
	l.addFields(event, keysAndValues...)
	event.Msg(msg)
}

func (l *zerologLogger) Warn(msg string, keysAndValues ...interface{}) {
	event := l.logger.Warn()
	l.addFields(event, keysAndValues...)
	event.Msg(msg)
}

func (l *zerologLogger) Fatal(msg string, keysAndValues ...interface{}) {
	event := l.logger.Fatal()
	l.addFields(event, keysAndValues...)
	event.Msg(msg)
}

func (l *zerologLogger) With(keysAndValues ...interface{}) Logger {
	ctx := l.logger.With()
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			key := keysAndValues[i].(string)
			value := keysAndValues[i+1]
			ctx = ctx.Interface(key, value)
		}
	}
	return &zerologLogger{logger: ctx.Logger()}
}

func (l *zerologLogger) addFields(event *zerolog.Event, keysAndValues ...interface{}) {
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			key := keysAndValues[i].(string)
			value := keysAndValues[i+1]

			switch v := value.(type) {
			case string:
				event.Str(key, v)
			case int:
				event.Int(key, v)
			case int64:
				event.Int64(key, v)
			case float64:
				event.Float64(key, v)
			case bool:
				event.Bool(key, v)
			case error:
				event.Err(v)
			case time.Duration:
				event.Dur(key, v)
			default:
				event.Interface(key, v)
			}
		}
	}
}

var global Logger

func init() {
	global = New("info")
}

func Info(msg string, keysAndValues ...interface{}) {
	global.Info(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	global.Error(msg, keysAndValues...)
}

func Debug(msg string, keysAndValues ...interface{}) {
	global.Debug(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...interface{}) {
	global.Warn(msg, keysAndValues...)
}

func Fatal(msg string, keysAndValues ...interface{}) {
	global.Fatal(msg, keysAndValues...)
}

// OBS: SetGlobal to define logger as global
func SetGlobal(l Logger) {
	global = l
}
