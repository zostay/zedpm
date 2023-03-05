package main

import (
	"fmt"
	"strings"
	"time"
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
		AddLines(newToBoundary - oldToBoundary)
	} else if newToBoundary < oldToBoundary {
		MoveUp(oldToBoundary - newToBoundary)
		ClearLines(oldToBoundary - newToBoundary)
		MoveUp(oldToBoundary - newToBoundary)
	}

	MoveUp(newToBoundary)

	for _, logLine := range s.log {
		ClearLine()
		fmt.Println(logLine)
	}

	WriteBoundary()

	for _, p := range s.progress {
		for _, l := range p.lines {
			ClearLine()
			fmt.Println(l)
		}
	}
}

func MoveUp(n int) {
	_, _ = fmt.Printf("\x1b[%dA", n)
}

func ClearLine() {
	_, _ = fmt.Print("\x1b[2K")
}

func ClearLines(n int) {
	_, _ = fmt.Print(strings.Repeat("\x1b[2K\n", n))
}

func AddLines(n int) {
	_, _ = fmt.Print(strings.Repeat("\n", n))
}

func WriteBoundary() {
	ClearLine()
	fmt.Println("---- ⏷ Active Tasks ⏷ ---- ⏶ Logs ⏶ ----")
}

func main() {
	var prev *state
	WriteBoundary()
	for i := range simulation {
		s := &simulation[i]
		s.Update(prev)
		prev = s
		time.Sleep(1 * time.Second)
	}
	MoveUp(1)
	ClearLine()
}
