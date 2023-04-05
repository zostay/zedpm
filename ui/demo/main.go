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

type change struct {
	delay  time.Duration
	op     ChangeOp
	widget int
	n      int
	line   string
}

var simChanges = []change{
	{0, OpAddWidget, 0, 6, ""},
	{0, OpSetWidget, 0, 0, "╔═══════╗    ┌──────┐    ┌───────┐    ┌─────┐"},
	{0, OpSetWidget, 0, 1, "║ Start ║ -> │ Mint │ -> │ Phase │ -> │ End │"},
	{0, OpSetWidget, 0, 2, "╚═══════╝    └──────┘    └───────┘    └─────┘"},
	{0, OpSetWidget, 0, 3, "╔════════╗"},
	{0, OpSetWidget, 0, 4, "║ Plugin ║"},
	{0, OpSetWidget, 0, 5, "╚════════╝"},
	{0, OpAddWidget, 0, 4, ""},
	{0, OpSetWidgetTitle, 1, 0, "---- Configuring plugins ---------------"},
	{100 * time.Millisecond, OpLog, 1, 0, "[Initialize] master: Configuring plugins..."},
	{100 * time.Millisecond, OpLog, 1, 0, "[Initialize] master: - Loading zedpm-plugin-changelog"},
	{100 * time.Millisecond, OpLog, 1, 0, "[Initialize] master: - Loading zedpm-plugin-git"},
	{100 * time.Millisecond, OpLog, 1, 0, "[Initialize] master: - Loading zedpm-plugin-github"},
	{100 * time.Millisecond, OpLog, 1, 0, "[Initialize] master: - Loading zedpm-plugin-goals"},
	{1 * time.Second, OpLog, 1, 0, "[Initialize] master: Complete."},
	{0, OpSetWidget, 0, 0, "┌───────┐    ╔══════╗    ┌───────┐    ┌─────┐"},
	{0, OpSetWidget, 0, 1, "│ Start │ -> ║ Mint ║ -> │ Phase │ -> │ End │"},
	{0, OpSetWidget, 0, 2, "└───────┘    ╚══════╝    └───────┘    └─────┘"},
	{0, OpSetWidget, 0, 3, "╔═══════╗    ┌───────┐    ┌───────┐    ┌─────┐    ┌─────┐    ┌────────┐    ┌──────────┐"},
	{0, OpSetWidget, 0, 4, "║ Setup ║ -> │ Check │ -> │ Begin │ -> │ Run │ -> │ End │ -> │ Finish │ -> │ Teardown │"},
	{0, OpSetWidget, 0, 5, "╚═══════╝    └───────┘    └───────┘    └─────┘    └─────┘    └────────┘    └──────────┘"},
	{0, OpDeleteWidget, 1, 0, ""},
	{0, OpAddWidget, 0, 4, ""},
	{0, OpSetWidgetTitle, 1, 0, "---- Changelog -------------------------"},
	{0, OpAddWidget, 0, 4, ""},
	{0, OpSetWidgetTitle, 2, 0, "---- Git -------------------------------"},
	{0, OpAddWidget, 0, 4, ""},
	{0, OpSetWidgetTitle, 3, 0, "---- Github ----------------------------"},
	{100 * time.Millisecond, OpSetWidget, 0, 3, "┌───────┐    ╔═══════╗    ┌───────┐    ┌─────┐    ┌─────┐    ┌────────┐    ┌──────────┐"},
	{0, OpSetWidget, 0, 4, "│ Setup │ -> ║ Check ║ -> │ Begin │ -> │ Run │ -> │ End │ -> │ Finish │ -> │ Teardown │"},
	{0, OpSetWidget, 0, 5, "└───────┘    ╚═══════╝    └───────┘    └─────┘    └─────┘    └────────┘    └──────────┘"},
	{0, OpLog, 1, 0, "[Check] zedpm-plugin-changelog: Linting changelog..."},
	{0, OpLog, 2, 0, "[Check] zedpm-plugin-git: Check worktree cleanliness..."},
	{800 * time.Millisecond, OpLog, 3, 0, "[Check] zedpm-plugin-github: ..."},
	{300 * time.Millisecond, OpLog, 1, 0, "[Check] zedpm-plugin-changelog: - Changes.md: PASS"},
	{1100 * time.Millisecond, OpLog, 1, 0, "[Check] zedpm-plugin-changelog: Complete."},
	{900 * time.Millisecond, OpLog, 2, 0, "[Check] zedpm-plugin-git: - Found HEAD"},
	{100 * time.Millisecond, OpLog, 2, 0, "[Check] zedpm-plugin-git: - HEAD branch matches expected target branch: master"},
	{100 * time.Millisecond, OpLog, 2, 0, "[Check] zedpm-plugin-git: - Listing remote references."},
	{200 * time.Millisecond, OpLog, 2, 0, "[Check] zedpm-plugin-git: - Local copy matches remote reference."},
	{500 * time.Millisecond, OpLog, 2, 0, "[Check] zedpm-plugin-git: - Local copy is clean."},
	{1200 * time.Millisecond, OpLog, 2, 0, "[Check] zedpm-plugin-git: - Work tree check: PASS"},
	{400 * time.Millisecond, OpLog, 2, 0, "[Check] zedpm-plugin-git: Complete."},
}

func main() {
	term := ui.NewTerminal(os.Stdout)
	state := ui.NewState(term, 4)
	for _, c := range simChanges {
		switch c.op {
		case OpAddWidget:
			state.AddWidget(ui.NewWidget(c.n))
		case OpDeleteWidget:
			state.DeleteWidget(c.widget)
		case OpLog:
			state.Log(c.widget, c.line)
		case OpSetWidget:
			state.Set(c.widget, c.n, c.line)
		case OpSetWidgetTitle:
			state.SetTitle(c.widget, c.line)
		}
		time.Sleep(c.delay)
	}
	state.Close()
}
