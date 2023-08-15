package log

// Interface is the logging interface that is provided to both the plugin
// logging system and the master logging system.
type Interface interface {
	Debug(line string, fields ...any)
	Info(line string, fields ...any)
	Warn(line string, fields ...any)
	Error(line string, fields ...any)
	Err(err error, fields ...any)

	With(args ...any) Interface

	StartAction(key, desc string, flags ...ActionFlag)
	TickAction(key string)
	MarkAction(key string, outcome Outcome)
}
