package ui

import (
	"github.com/hashicorp/go-hclog"
)

// ProgressAdapter sends logs to the progress.
type ProgressAdapter struct {
	progress *Progress

	// because hclog doesn't seem to pay attention to this!?
	minLevel hclog.Level
}

func NewSinkAdapter(progress *Progress, minLevel hclog.Level) *ProgressAdapter {
	return &ProgressAdapter{progress, minLevel}
}

func (p *ProgressAdapter) Accept(
	name string,
	level hclog.Level,
	msg string,
	args ...any,
) {
	if level >= p.minLevel {
		p.progress.Log(name, level.String(), msg, args...)
	}
}
