package log

import (
	"github.com/hashicorp/go-hclog"
)

// Logger is the internal logger that is aware of the state of the process as it
// moves from one stage to the next. Plugins should use
// github.com/zostay/zedpm/pkg/log.Logger instead.
type Logger struct {
	hclog.Logger
	tnReq  chan struct{}
	tnChan chan string
}

func (l *Logger) taskName() string {
	l.tnReq <- struct{}{}
	return <-l.tnChan
}

func (l *Logger) argsWithTaskName(args []any) []any {
	taskName := l.taskName()
	if taskName != "" {
		return append(args, "task", l.taskName())
	}
	return args
}

func NewPluginLogWrapper(logger hclog.Logger) (hclog.Logger, chan<- string) {
	l := &Logger{
		Logger: logger,
		tnReq:  make(chan struct{}),
		tnChan: make(chan string),
	}

	tnXfer := make(chan string)

	go func() {
		taskName := ""
		for {
			select {
			case tn, stillOpen := <-tnXfer:
				if !stillOpen {
					return
				}
				taskName = tn
			case <-l.tnReq:
				l.tnChan <- taskName
			}
		}
	}()

	return l, tnXfer
}

func (l *Logger) Log(level hclog.Level, msg string, args ...any) {
	l.Logger.Log(level, msg, l.argsWithTaskName(args)...)
}

func (l *Logger) Trace(msg string, args ...any) {
	l.Logger.Trace(msg, l.argsWithTaskName(args)...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.Logger.Debug(msg, l.argsWithTaskName(args)...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.Logger.Info(msg, l.argsWithTaskName(args)...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.Logger.Warn(msg, l.argsWithTaskName(args)...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.Logger.Error(msg, l.argsWithTaskName(args)...)
}

func (l *Logger) ImpliedArgs() []any {
	return l.argsWithTaskName(l.Logger.ImpliedArgs())
}

func (l *Logger) With(args ...any) hclog.Logger {
	return &Logger{
		Logger: l.Logger.With(args...),
		tnReq:  l.tnReq,
		tnChan: l.tnChan,
	}
}

func (l *Logger) Named(name string) hclog.Logger {
	return &Logger{
		Logger: l.Logger.Named(name),
		tnReq:  l.tnReq,
		tnChan: l.tnChan,
	}
}

func (l *Logger) ResetNamed(name string) hclog.Logger {
	return &Logger{
		Logger: l.Logger.ResetNamed(name),
		tnReq:  l.tnReq,
		tnChan: l.tnChan,
	}
}
