package main

import (
	"os"
	"time"

	"github.com/zostay/zedpm/ui"
)

type ChangeOp int

const (
	OpLog ChangeOp = iota + 1
	OpAddWidget
	OpDeleteWidget
	OpSetWidget
	OpSetWidgetTitle
)

type MyID int

const (
	ProgressBar MyID = iota + 1
	InitializeLog
	ChangelogLog
	GitLog
	GithubLog
)

var wm = map[MyID]ui.WidgetID{
	ProgressBar:   0,
	InitializeLog: 0,
	ChangelogLog:  0,
	GitLog:        0,
	GithubLog:     0,
}

type change struct {
	delay  time.Duration
	op     ChangeOp
	widget MyID
	n      int
	line   string
}

var (
	red    = "\U0001f534"
	yellow = "\U0001f7e1"
	green  = "\U0001f7e2"
)

var simChanges = []change{
	{0, OpAddWidget, ProgressBar, 4, ""},
	{0, OpSetWidget, ProgressBar, 0, yellow + " Initialize"},
	{0, OpSetWidget, ProgressBar, 1, red + " Mint"},
	{0, OpSetWidget, ProgressBar, 2, red + " Phase"},
	{0, OpSetWidget, ProgressBar, 3, red + " Quit"},
	{0, OpAddWidget, InitializeLog, 4, ""},
	{0, OpSetWidgetTitle, InitializeLog, 0, "Configuring plugins"},
	{100 * time.Millisecond, OpLog, InitializeLog, 0, "[Initialize] master: Configuring plugins..."},
	{100 * time.Millisecond, OpLog, InitializeLog, 0, "[Initialize] master: - Loading zedpm-plugin-changelog"},
	{100 * time.Millisecond, OpLog, InitializeLog, 0, "[Initialize] master: - Loading zedpm-plugin-git"},
	{100 * time.Millisecond, OpLog, InitializeLog, 0, "[Initialize] master: - Loading zedpm-plugin-github"},
	{100 * time.Millisecond, OpLog, InitializeLog, 0, "[Initialize] master: - Loading zedpm-plugin-goals"},
	{1 * time.Second, OpLog, InitializeLog, 0, "[Initialize] master: Complete."},
	{0, OpSetWidget, ProgressBar, 0, green + " Initialize"},
	{0, OpSetWidget, ProgressBar, 1, yellow + " Mint [Setup]"},
	{0, OpDeleteWidget, InitializeLog, 0, ""},
	{0, OpAddWidget, ChangelogLog, 4, ""},
	{0, OpSetWidgetTitle, ChangelogLog, 0, "Changelog"},
	{0, OpAddWidget, GitLog, 4, ""},
	{0, OpSetWidgetTitle, GitLog, 0, "Git"},
	{0, OpAddWidget, GithubLog, 4, ""},
	{0, OpSetWidgetTitle, GithubLog, 0, "Github"},
	{0 * time.Millisecond, OpSetWidget, ProgressBar, 1, yellow + " Mint [Check]"},
	{100, OpLog, ChangelogLog, 0, "[Check] zedpm-plugin-changelog: Linting changelog..."},
	{0, OpLog, GitLog, 0, "[Check] zedpm-plugin-git: Check worktree cleanliness..."},
	{800 * time.Millisecond, OpLog, GithubLog, 0, "[Check] zedpm-plugin-github: ..."},
	{300 * time.Millisecond, OpLog, ChangelogLog, 0, "[Check] zedpm-plugin-changelog: - Changes.md: PASS"},
	{900 * time.Millisecond, OpLog, GitLog, 0, "[Check] zedpm-plugin-git: - Found HEAD"},
	{100 * time.Millisecond, OpLog, GitLog, 0, "[Check] zedpm-plugin-git: - HEAD branch matches expected target branch: master"},
	{100 * time.Millisecond, OpLog, GitLog, 0, "[Check] zedpm-plugin-git: - Listing remote references."},
	{200 * time.Millisecond, OpLog, GitLog, 0, "[Check] zedpm-plugin-git: - Local copy matches remote reference."},
	{500 * time.Millisecond, OpLog, GitLog, 0, "[Check] zedpm-plugin-git: - Local copy is clean."},
	{1200 * time.Millisecond, OpLog, GitLog, 0, "[Check] zedpm-plugin-git: - Work tree check: PASS"},
	{0, OpSetWidget, ProgressBar, 1, yellow + " Mint [Run:30]"},
	{300 * time.Millisecond, OpLog, GitLog, 0, "[Run:30] zedpm-plugin-git: - Created git branch for managing the release"},
	{0, OpSetWidget, ProgressBar, 1, yellow + " Mint [Run:50]"},
	{200 * time.Millisecond, OpLog, ChangelogLog, 0, "[Run:50] zedpm-plugin-changelog: - Applied changes to changelog to fixup for release"},
	{100 * time.Millisecond, OpLog, ChangelogLog, 0, "[Run:55] zedpm-plugin-changelog: - Changelog linted for release: PASS"},
	{0, OpSetWidget, ProgressBar, 1, yellow + " Mint [Run:55]"},
	{1100 * time.Millisecond, OpLog, ChangelogLog, 0, "[Complete] zedpm-plugin-changelog: Mint Phase Complete"},
	{0, OpSetWidget, ProgressBar, 1, yellow + " Mint [End:70]"},
	{300 * time.Millisecond, OpLog, GitLog, 0, "[End:70] zedpm-plugin-git: - Added Files and committing changes to git"},
	{0, OpSetWidget, ProgressBar, 1, yellow + " Mint [End:75]"},
	{700 * time.Millisecond, OpLog, GitLog, 0, "[End:75] zedpm-plugin-git: - Pushed release branch to remote repository"},
	{400 * time.Millisecond, OpLog, GitLog, 0, "[Complete] zedpm-plugin-git: Mint Phase Complete"},
	{1400 * time.Millisecond, OpLog, GithubLog, 0, "[End:80] zedpm-plugin-github: - Created Github pull request"},
	{600 * time.Millisecond, OpLog, GithubLog, 0, "[Complete] zedpm-plugin-github: Mint Phase Complete"},
	{0, OpSetWidget, ProgressBar, 1, green + " Mint"},
	{0, OpSetWidget, ProgressBar, 2, yellow + " Publish"},
}

func main() {
	term := ui.NewTerminal(os.Stdout)
	state := ui.NewState(term, 4)
	for _, c := range simChanges {
		switch c.op {
		case OpAddWidget:
			wm[c.widget] = state.AddWidget(c.n)
		case OpDeleteWidget:
			state.DeleteWidget(wm[c.widget])
		case OpLog:
			state.Log(wm[c.widget], c.line)
		case OpSetWidget:
			state.Set(wm[c.widget], c.n, c.line)
		case OpSetWidgetTitle:
			state.SetTitle(wm[c.widget], c.line)
		}
		time.Sleep(c.delay)
	}
	time.Sleep(2 * time.Second)
	state.Close()
}
