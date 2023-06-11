package ui

import "github.com/zostay/zedpm/pkg/log"

type statusIcon string

const (
	redCircle    statusIcon = "\U0001f534"
	yellowCircle statusIcon = "\U0001f7e1"
	greenCircle  statusIcon = "\U0001f7e2"
	purpleCircle statusIcon = "\U0001f3e3"
)

type widgetLine interface {
	SetActionKey(action string)
	ActionKey() string
	String() string
}

type widgetOutcomeLine interface {
	Outcome() log.Outcome
	SetOutcome(outcome log.Outcome)
}

type widgetFlagLine interface {
	AddFlags(flags ...string)
	RemoveFlags(flags ...string)
}

type widgetTickLine interface {
	IncTick()
}

type widgetTitleLine interface {
	Title() string
	SetTitle(title string)
}

type widgetStatusLine interface {
	Icon() statusIcon
	SetIcon(icon statusIcon)
}
