package log

import (
	"fmt"
	"io"

	"github.com/zostay/go-std/slices"
)

// Outcome is used to record the outcome of some running action.
type Outcome string

const (
	Working Outcome = ""
	Fail    Outcome = "fail"
	Pass    Outcome = "pass"
	Skip    Outcome = "skip"
	Retry   Outcome = "retry"
)

// IsFinal returns true if the Outcome is considered a terminal state.
func (o Outcome) IsFinal() bool {
	return o == Fail || o == Pass || o == Skip
}

// ActionFlag provides additional flags for actions
type ActionFlag string

const (
	// None is not a flag.
	None ActionFlag = "none"

	// Spinning requests a manually ticked spinner.
	Spinning ActionFlag = "spin"

	// AutoSpinning requests an automatically ticked spinner.
	AutoSpinning ActionFlag = "auto-spin"

	// Plain requests that the log line be printed without spinners, outcomes,
	// etc. Just the message itself.
	Plain ActionFlag = "plain"
)

type LowLevelLogger interface {
	Debug(line string, fields ...any)
	Info(line string, fields ...any)
	Warn(line string, fields ...any)
	Error(line string, fields ...any)
}

// Logger is provided to plugins to allow them to interact with the parent task
// in a standard way. Other language implementations can mimic the behavior here
// to achieve the same results. Under the hood is just a structured logging
// system passing information back to the master process via tagged values.
type Logger struct {
	fields  []any
	log     LowLevelLogger
	actions map[string]string
}

// New creates a new logger for use with a plugin.
func New(log LowLevelLogger) *Logger {
	return &Logger{
		fields:  []any{},
		actions: make(map[string]string, 10),
		log:     log,
	}
}

// allFields just glues the given args to the With fields.
func (l *Logger) allFields(argLists ...[]any) []any {
	argLists = slices.Unshift(argLists, l.fields)
	return slices.Concat(argLists...)
}

// StartAction logs the given message, but marks the action as being a progress
// item that has a followup. You may pass zero or more ActionFlags to mark this
// action as having additional features. Currently, flags include:
//
//	log.None - no-op flag
//	log.Spin - request a manually ticked spinner
//	log.AutoSpin - request an automatically ticked spinner
//	log.Plain - just the message, no spinners, outcomes, etc.
func (l *Logger) StartAction(key, desc string, flags ...ActionFlag) {
	l.actions[key] = desc
	l.Info(desc, "@action", key, "@actionFlags", flags)
}

// TickAction logs the original message again, but is intended to be used to
// make a spinner move a tick.
func (l *Logger) TickAction(key string) {
	desc := l.actions[key]
	l.Debug(desc, "@action", key, "@tick", 1)
}

// MarkAction logs the original message with the given outcome.
func (l *Logger) MarkAction(key string, outcome Outcome) {
	desc := l.actions[key]
	delete(l.actions, key)
	l.Info(desc+": %[@outcome]s", "@action", key, "@outcome", outcome)
}

// Info logs the given info log.
func (l *Logger) Info(fmt string, args ...any) {
	args = l.allFields(args)
	l.log.Info(Smprintf(fmt, args...), args...)
}

// Error logs the given error log.
func (l *Logger) Error(fmt string, args ...any) {
	args = l.allFields(args)
	l.log.Error(Smprintf(fmt, args...), args...)
}

// Err logs the given error as an error log. It will add any details attached to
// the log (via the log.WithDetails interface) as additional arguments.
func (l *Logger) Err(err error, args ...any) {
	details := GetDetails(err)
	args = l.allFields(details, args)
	args = append(args, "error", err)
	l.log.Error(fmt.Sprintf("error: %v", err), args...)
}

// Warn logs the given warning log.
func (l *Logger) Warn(fmt string, args ...any) {
	args = l.allFields(args)
	l.log.Warn(Smprintf(fmt, args...), args...)
}

// Debug logs the given debug log.
func (l *Logger) Debug(fmt string, args ...any) {
	args = l.allFields(args)
	l.log.Debug(Smprintf(fmt, args...), args...)
}

type logWriter struct {
	logFunc func(string, ...any)
}

func (w *logWriter) Write(p []byte) (int, error) {
	n := len(p)
	w.logFunc(string(p))
	return n, nil
}

// Output returns an io.Writer that can be used to write to the given log level.
func (l *Logger) Output(level Level) io.Writer {
	var logFunc func(string, ...any)
	switch level {
	case LevelDebug:
		logFunc = l.Debug
	case LevelInfo:
		logFunc = l.Info
	case LevelWarn:
		logFunc = l.Warn
	case LevelError:
		logFunc = l.Error
	default:
		panic("unknown level")
	}

	return &logWriter{logFunc}
}

// With returns a new logger with the given args included with every log sent.
func (l *Logger) With(args ...any) Interface {
	fields := l.allFields(args)
	return &Logger{
		fields:  fields,
		actions: l.actions,
		log:     l.log,
	}
}
