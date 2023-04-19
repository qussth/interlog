package interlog

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
)

type Reported struct {
}

// logger struct
type logger struct {
	zero zerolog.Logger
}

// Logger ...
type Logger interface {
	Debug(message string, values []Value)
	Info(message string, values []Value)
	InfoToSentry(message string, values []Value)
	Warn(message string, values []Value)
	Error(err error, values []Value)
	Panic(err error, values []Value)
}

type Values []Value

type Value struct {
	Key     string
	Payload interface{}
}

var Message = zerolog.MessageFieldName

// New function
func New() Logger {
	l := &logger{
		zero: zerolog.New(zerolog.ConsoleWriter{
			Out: os.Stdout,
			FormatLevel: func(i interface{}) string {
				return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
			},
			FormatTimestamp: func(i interface{}) string {
				return time.Now().Format("2006-01-02 15:04:05")
			},
		}),
	}

	level := strings.ToLower(os.Getenv("MODE"))
	switch level {
	case "debug":
		l.zero = l.zero.Level(zerolog.DebugLevel)
	case "warn":
		l.zero = l.zero.Level(zerolog.WarnLevel)
	case "error":
		l.zero = l.zero.Level(zerolog.ErrorLevel)
	case "panic":
		l.zero = l.zero.Level(zerolog.PanicLevel)
	default:
		l.zero = l.zero.Level(zerolog.InfoLevel)
	}

	return l
}

func (l *logger) InitializeSentry(sentryOptions sentry.ClientOptions) error {
	return sentry.Init(sentryOptions)
}

// Debug func
func (l *logger) Debug(message string, values []Value) {
	event := l.zero.Debug()

	l.appendInterfaces(event, values)
	event.Msg(message)

}

// Info func
func (l *logger) Info(message string, values []Value) {
	event := l.zero.Info()

	l.appendInterfaces(event, values)
	event.Msg(message)
}

// InfoToSentry func
// Send info message also in Sentry
func (l *logger) InfoToSentry(message string, values []Value) {
	event := l.zero.Info()

	l.appendInterfaces(event, values)
	event.Msg(message)
}

// Warn func
func (l *logger) Warn(message string, values []Value) {
	event := l.zero.Warn()

	l.appendInterfaces(event, values)
	event.Msg(message)
}

// Error func
// pass `zerolog.MessageFieldName` field in values to set Msg
func (l *logger) Error(err error, values []Value) {
	event := l.zero.Error().Err(err)

	event.Msg(l.iface(values, event))
	sentry.CaptureException(err)
}

// Panic func
// will invoke panic with err.Error()
// pass zerolog.MessageFieldName field in values to set Msg
func (l *logger) Panic(err error, values []Value) {
	event := l.zero.Panic().Err(err)

	event.Msg(l.iface(values, event))
	sentry.CaptureException(err)
}

func (l *logger) iface(values []Value, event *zerolog.Event) string {
	var msgPassed string

	valIndex := 1

	for _, value := range values {
		if value.Key == zerolog.MessageFieldName {
			msgPassed = fmt.Sprintf("%v", value.Payload)
			continue
		}

		event = event.Interface(value.Key, fmt.Sprintf("%d:%v", valIndex, value.Payload))
		valIndex++
	}

	return msgPassed
}

func (l *logger) appendInterfaces(event *zerolog.Event, values []Value) *zerolog.Event {
	for i, value := range values {
		event = event.Interface(fmt.Sprintf("%d:%s", i+1, value.Key), value.Payload)
	}

	return event
}
