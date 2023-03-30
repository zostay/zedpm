package main

import (
	"fmt"
	"time"

	"github.com/zostay/zedpm/ui"
)

type progress struct {
	lines []string
}

type state struct {
	log      []string
	progress []progress
}

var simulation = []state{
	{
		log: []string{
			"[Initialize] master: Configuring plugins...",
		},
		progress: []progress{
			{
				lines: []string{
					"[Initialize] master: Configuring plugins...",
				},
			},
		},
	},
	{
		log: []string{
			"[Initialize] master:  - Loading zedpm-plugin-changelog",
		},
		progress: []progress{
			{
				lines: []string{
					"[Initialize] master: Configuring plugins...",
					" - Loading zedpm-plugin-changelog",
					"",
					"",
				},
			},
		},
	},
	{
		log: []string{
			"[Initialize] master:  - Loading zedpm-plugin-git",
		},
		progress: []progress{
			{
				lines: []string{
					"[Initialize] master: Configuring plugins...",
					" - Loading zedpm-plugin-changelog",
					" - Loading zedpm-plugin-git",
					"",
				},
			},
		},
	},
	{
		log: []string{
			"[Initialize] master:  - Loading zedpm-plugin-github",
		},
		progress: []progress{
			{
				lines: []string{
					"[Initialize] master: Configuring plugins...",
					" - Loading zedpm-plugin-changelog",
					" - Loading zedpm-plugin-git",
					" - Loading zedpm-plugin-github",
				},
			},
		},
	},
	{
		log: []string{
			"[Initialize] master:  - Loading zedpm-plugin-goals",
		},
		progress: []progress{
			{
				lines: []string{
					"[Initialize] Configuring plugins...",
					" - Loading zedpm-plugin-git",
					" - Loading zedpm-plugin-github",
					" - Loading zedpm-plugin-goals",
				},
			},
		},
	},
	{
		log: []string{
			"[Initialize] master: Complete.",
		},
	},
	{},
	{
		log: []string{
			"[Check] zedpm-plugin-changelog: Linting changelog...",
			"[Check] zedpm-plugin-git: Check worktree cleanliness...",
			"[Check] zedpm-plugin-github: ...",
		},
		progress: []progress{
			{
				lines: []string{
					"[Check] zedpm-plugin-changelog: Linting changelog...",
					"",
					"",
					"",
				},
			},
			{
				lines: []string{
					"[Check] zedpm-plugin-git: Check worktree cleanliness...",
					"",
					"",
					"",
				},
			},
			{
				lines: []string{
					"[Check] zedpm-plugin-github: ...",
					"",
					"",
					"",
				},
			},
		},
	},
	{
		log: []string{
			"[Check] zedpm-plugin-changelog:  - Changes.md: PASS",
			"[Check] zedpm-plugin-changelog: Complete.",
			"[Check] zedpm-plugin-git:  - Found HEAD",
		},
		progress: []progress{
			{
				lines: []string{
					"[Check] zedpm-plugin-changelog: Complete.",
					" - Changes.md: PASS",
					"",
					"",
				},
			},
			{
				lines: []string{
					"[Check] zedpm-plugin-git: Check worktree cleanliness...",
					" - Found HEAD",
					"",
					"",
				},
			},
			{
				lines: []string{
					"[Check] zedpm-plugin-github: ...",
					"",
					"",
					"",
				},
			},
		},
	},
	{
		log: []string{
			"[Check] zedpm-plugin-git:  - HEAD branch matches expected target branch: master",
		},
		progress: []progress{
			{
				lines: []string{
					"[Check] zedpm-plugin-changelog: Complete.",
					" - Changes.md: PASS",
					"",
					"",
				},
			},
			{
				lines: []string{
					"[Check] zedpm-plugin-git: Check worktree cleanliness...",
					" - Found HEAD",
					" - HEAD branch matches expected target branch: master",
					"",
				},
			},
			{
				lines: []string{
					"[Check] zedpm-plugin-github: ...",
					"",
					"",
					"",
				},
			},
		},
	},
	{
		log: []string{
			"[Check] zedpm-plugin-git:  - Listing remote references.",
		},
		progress: []progress{
			{
				lines: []string{
					"[Check] zedpm-plugin-changelog: Complete.",
					" - Changes.md: PASS",
					"",
					"",
				},
			},
			{
				lines: []string{
					"[Check] zedpm-plugin-git: Check worktree cleanliness...",
					" - Found HEAD",
					" - HEAD branch matches expected target branch: master",
					" - Listing remote references.",
				},
			},
			{
				lines: []string{
					"[Check] zedpm-plugin-github: ...",
					"",
					"",
					"",
				},
			},
		},
	},
	{
		log: []string{
			"[Check] zedpm-plugin-git:  - Local copy matches remote reference.",
		},
		progress: []progress{
			{
				lines: []string{
					"[Check] zedpm-plugin-changelog: Complete.",
					" - Changes.md: PASS",
					"",
					"",
				},
			},
			{
				lines: []string{
					"[Check] zedpm-plugin-git: Check worktree cleanliness...",
					" - HEAD branch matches expected target branch: master",
					" - Listing remote references.",
					" - Local copy matches remote reference.",
				},
			},
			{
				lines: []string{
					"[Check] zedpm-plugin-github: ...",
					"",
					"",
					"",
				},
			},
		},
	},
	{
		log: []string{
			"[Check] zedpm-plugin-git:  - Local copy is clean.",
		},
		progress: []progress{
			{
				lines: []string{
					"[Check] zedpm-plugin-changelog: Complete.",
					" - Changes.md: PASS",
					"",
					"",
				},
			},
			{
				lines: []string{
					"[Check] zedpm-plugin-git: Check worktree cleanliness...",
					" - Listing remote references.",
					" - Local copy matches remote reference.",
					" - Local copy is clean.",
				},
			},
			{
				lines: []string{
					"[Check] zedpm-plugin-github: ...",
					"",
					"",
					"",
				},
			},
		},
	},
	{
		log: []string{
			"[Check] zedpm-plugin-git:  - Work tree check: PASS",
			"[Check] zedpm-plugin-git: Complete.",
		},
		progress: []progress{
			{
				lines: []string{
					"[Check] zedpm-plugin-changelog: Complete.",
					" - Changes.md: PASS",
					"",
					"",
				},
			},
			{
				lines: []string{
					"[Check] zedpm-plugin-git: Complete.",
					" - Local copy matches remote reference.",
					" - Local copy is clean.",
					" - Work tree check: PASS",
				},
			},
			{
				lines: []string{
					"[Check] zedpm-plugin-github: ...",
					"",
					"",
					"",
				},
			},
		},
	},
	{},
}

const (
	OpLog = iota + 1
	OpAddWidget
	OpDeleteWidget
)

type change struct {
	delay  time.Duration
	op     int
	widget int
	line   string
}

var simChanges = []change{
	{1 * time.Second, OpAddWidget, 4, ""},
	{100 * time.Millisecond, OpLog, 0, "[Initialize] master: Configuring plugins..."},
	{100 * time.Millisecond, OpLog, 0, "[Initialize] master: - Loading zedpm-plugin-changelog"},
	{100 * time.Millisecond, OpLog, 0, "[Initialize] master: - Loading zedpm-plugin-git"},
	{100 * time.Millisecond, OpLog, 0, "[Initialize] master: - Loading zedpm-plugin-github"},
	{100 * time.Millisecond, OpLog, 0, "[Initialize] master: - Loading zedpm-plugin-goals"},
	{1 * time.Second, OpLog, 0, "[Initialize] master: Complete."},
	{0, OpDeleteWidget, 0, ""},
	{0, OpAddWidget, 4, ""},
	{0, OpLog, 0, "[Check] zedpm-plugin-changelog: Linting changelog..."},
	{0, OpAddWidget, 4, ""},
	{0, OpLog, 1, "[Check] zedpm-plugin-git: Check worktree cleanliness..."},
	{0, OpAddWidget, 4, ""},
	{800 * time.Millisecond, OpLog, 2, "[Check] zedpm-plugin-github: ..."},
	{300 * time.Millisecond, OpLog, 0, "[Check] zedpm-plugin-changelog: - Changes.md: PASS"},
	{1100 * time.Millisecond, OpLog, 0, "[Check] zedpm-plugin-changelog: Complete."},
	{900 * time.Millisecond, OpLog, 1, "[Check] zedpm-plugin-git: - Found HEAD"},
	{100 * time.Millisecond, OpLog, 1, "[Check] zedpm-plugin-git: - HEAD branch matches expected target branch: master"},
	{100 * time.Millisecond, OpLog, 1, "[Check] zedpm-plugin-git: - Listing remote references."},
	{200 * time.Millisecond, OpLog, 1, "[Check] zedpm-plugin-git: - Local copy matches remote reference."},
	{500 * time.Millisecond, OpLog, 1, "[Check] zedpm-plugin-git: - Local copy is clean."},
	{1200 * time.Millisecond, OpLog, 1, "[Check] zedpm-plugin-git: - Work tree check: PASS"},
	{400 * time.Millisecond, OpLog, 1, "[Check] zedpm-plugin-git: Complete."},
}

func (s *state) MovementsToBoundary() int {
	l := 1
	for _, p := range s.progress {
		l += len(p.lines)
	}
	return l
}

func (s *state) Update(prev *state) {
	if prev == nil {
		prev = &state{}
	}

	var (
		newToBoundary = s.MovementsToBoundary()
		oldToBoundary = prev.MovementsToBoundary()
	)

	if newToBoundary > oldToBoundary {
		ui.AddLines(newToBoundary - oldToBoundary)
	} else if newToBoundary < oldToBoundary {
		ui.MoveUp(oldToBoundary - newToBoundary)
		ui.ClearLines(oldToBoundary - newToBoundary)
		ui.MoveUp(oldToBoundary - newToBoundary)
	}

	ui.MoveUp(newToBoundary)

	for _, logLine := range s.log {
		ui.ClearLine()
		fmt.Println(logLine)
	}

	ui.WriteBoundary()

	for _, p := range s.progress {
		for _, l := range p.lines {
			ui.ClearLine()
			fmt.Println(l)
		}
	}
}

func origMain() {
	var prev *state
	ui.WriteBoundary()
	for i := range simulation {
		s := &simulation[i]
		s.Update(prev)
		prev = s
		time.Sleep(1 * time.Second)
	}
	ui.MoveUp(1)
	ui.ClearLine()
}

func main() {
	state := ui.NewState(3)
	for _, c := range simChanges {
		switch c.op {
		case OpAddWidget:
			state = state.AddWidget(ui.NewWidget(c.widget))
		case OpDeleteWidget:
			state = state.DeleteWidget(c.widget)
		case OpLog:
			state.Log(c.widget, c.line)
		}
		time.Sleep(c.delay)
	}
	state.Close()
}
