package ui

import (
	"bytes"
)

type ProgressWriter struct {
	name     string
	level    string
	progress *Progress
	buffer   bytes.Buffer
}

func NewWriter(name, level string, progress *Progress) *ProgressWriter {
	return &ProgressWriter{
		name:     name,
		level:    level,
		progress: progress,
	}
}

func (w *ProgressWriter) Write(p []byte) (int, error) {
	n, _ := w.buffer.Write(p)

	for {
		line, err := w.buffer.ReadBytes('\n')
		if err != nil {
			return n, nil
		}

		w.progress.Log(w.name, w.level, string(line))
	}
}
