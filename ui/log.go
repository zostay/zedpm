package ui

import (
	"github.com/hashicorp/go-hclog"
)

// ProgressAdapter sends logs to the progress.
type ProgressAdapter struct {
	progress *Progress
}

func NewSinkAdapter(progress *Progress) *ProgressAdapter {
	return &ProgressAdapter{progress}
}

func (p *ProgressAdapter) Accept(
	name string,
	level hclog.Level,
	msg string,
	args ...any,
) {
	p.progress.Log(name, level.String(), msg, args...)
}
