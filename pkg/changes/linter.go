package changes

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/coreos/go-semver/semver"
)

// CheckMode identifies the kind of linting to perform on a changelog.
type CheckMode int

const (
	// CheckStandard merely checks that the syntax of the changelog file appears
	// to be sane.
	CheckStandard CheckMode = 0 + iota

	// CheckPreRelease checks the syntax and then checks that the first section
	// is a WIP section.
	CheckPreRelease

	// CheckRelease checks the syntax and then checks that the first section is
	// not a WIP section.
	CheckRelease
)

// Linter is the object with methods for performing changelog linting. The
// expected format for changelogs being linted by this method looks like this:
//
//	WIP  TBD
//
//	 * This is a change line for the work in progress.
//	 * This is a long line that will end up wrapped in the file because it is
//	   long and long lines are tedious if you have to scroll horizontally.
//
//	v1.0  2023-03-04
//
//	 * This is the change for the latest release.
//
// A changelog is made up of one or more sections. Each section starts with a
// heading of the form "<vstring> <date>" with exactly two spaces between the
// vstring and date. The first section may be "WIP TBD". (During a pre-release
// check, the first section must be a WIP section and during a release check,
// the first section must not be a WIP section.) This heading must be followed
// by a blank line. Then there are one or more bullet lines, each of which must
// start with " * " at the very start of the line. A bullet line may be
// continued on one or more following lines by ensuring that at least 3 spaces
// start each of those lines. Finally, prior to starting a new section head
// there must be another blank. Other than that, anything goes.
type Linter struct {
	r    io.Reader
	mode CheckMode
}

// Failure describes a line that has some defect and the description of the
// defect.
type Failure struct {
	Line    int
	Message string
}

// Failures is a list of changelog errors.
type Failures []Failure

// String outputs the Failures as a bulleted list.
func (fs Failures) String() string {
	buf := &strings.Builder{}
	for i, f := range fs {
		if i > 0 {
			_, _ = fmt.Fprint(buf, "\n")
		}
		_, _ = fmt.Fprintf(buf, " * Line %d: %s", f.Line, f.Message)
	}
	return buf.String()
}

// Error is an error made up of Failures and is returned by the Linter.Check
// method when one or more problems are detected in the changelog.
type Error struct {
	Failures
}

// Error returns the error message as a string.
func (e *Error) Error() string {
	return fmt.Sprintf("Change log linter check failed:\n%s", e.Failures.String())
}

// NewLinter constructs a Linter for checking a changelog file.
func NewLinter(r io.Reader, mode CheckMode) *Linter {
	return &Linter{r, mode}
}

// checkStatus is an internal structure used to track the current state of
// linter as each line is processed.
type checkStatus struct {
	previousVersion *semver.Version // the semantic version last read from a section head.
	previousDate    string          // the date last read from a section head.
	previousLine    int             // the line number of the last read section head.

	previousLineWasBlank  bool // true when the current line follows a blank line
	previousLineWasBullet bool // true when the current line follows a bullet line (start or continuation)

	Failures // the accumulated list of failures
}

// fail adds a new failure to the checkStatus.
func (s *checkStatus) fail(
	lineNumber int,
	msg string,
) {
	if s.Failures == nil {
		s.Failures = Failures{}
	}
	s.Failures = append(s.Failures, Failure{lineNumber, msg})
}

// failf adds a new failure to the checkStatus with printf formatting.
func (s *checkStatus) failf(
	lineNumber int,
	f string,
	args ...any,
) {
	s.fail(lineNumber, fmt.Sprintf(f, args...))
}

// Check executes the changelog linter against the reader. It scans each line of
// the changelog and checks it for problems. If there are no errors, this method
// returns nil. If one or more problems are detected, they will be returned as
// an Error.
func (l *Linter) Check() error {
	status := checkStatus{}

	s := bufio.NewScanner(l.r)
	n := 0
	for s.Scan() {
		n++
		l.checkLine(n, s.Text(), &status)
	}

	if len(status.Failures) > 0 {
		return &Error{status.Failures}
	} else {
		return nil
	}
}

var (
	versionHeading      = regexp.MustCompile(`^v(\d\S+) {2}(20\d\d-\d\d-\d\d)$`) // "vstring  date" lines
	logLineStart        = regexp.MustCompile(`^ \* (.*)$`)                       // bullet start lines
	logLineContinuation = regexp.MustCompile(`^ {3}(.*)$`)                       // bullet continuation lines
	blankLine           = regexp.MustCompile(`^$`)                               // completely empty lines
	whitespaceLine      = regexp.MustCompile(`^\s+$`)                            // lines containing whitespace
)

// checkLine is the workhorse function that checks that each line makes sense
// given the current status of the checks up to this point.
func (l *Linter) checkLine(
	lineNumber int,
	line string,
	status *checkStatus,
) {
	lineIsBlank := false
	lineIsBullet := false

	defer func() {
		status.previousLineWasBlank = lineIsBlank
		status.previousLineWasBullet = lineIsBullet
	}()

	if line == "WIP" || line == "WIP  TBD" {
		if lineNumber > 1 {
			status.fail(lineNumber, "WIP found after line 1")
		}

		if l.mode == CheckRelease {
			status.fail(lineNumber, "Found WIP line during release")
		}

		status.previousLine = lineNumber

		return
	}

	// we shouldn't get here if the first line is a WIP line
	if l.mode == CheckPreRelease && lineNumber == 1 {
		status.fail(lineNumber, "WIP not found during pre-release check")
	}

	if m := versionHeading.FindStringSubmatch(line); m != nil {
		ver, date := m[1], m[2]
		version, err := semver.NewVersion(ver)
		if err != nil {
			status.fail(lineNumber, "Unable to parse version number in heading")

			// this is fatal for this line, checks cannot continue
			status.previousLine = lineNumber
		}

		// version and date are in descending order in a changelog

		if status.previousVersion != nil && status.previousVersion.LessThan(*version) {
			status.failf(lineNumber, "version error %s < %s from line %d",
				version, status.previousVersion, status.previousLine)
		}

		if status.previousDate != "" && status.previousDate < date {
			status.failf(lineNumber, "date error %s < %s from line %d",
				date, status.previousDate, status.previousLine)
		}

		if lineNumber != 1 && !status.previousLineWasBlank {
			status.fail(lineNumber, "version heading line missing blank line before it")
		}

		status.previousVersion = version
		status.previousDate = date
		status.previousLine = lineNumber

		return
	}

	if m := logLineStart.FindStringSubmatch(line); m != nil {
		if status.previousLine == 0 {
			status.fail(lineNumber, "log bullet before first version heading or WIP")
		}

		if status.previousLine > 0 && lineNumber-1 == status.previousLine {
			status.fail(lineNumber, "missing blank line before log bullet")
		}

		if status.previousLine > 0 && lineNumber > status.previousLine+2 && status.previousLineWasBlank {
			status.fail(lineNumber, "extra blank line before log bullet")
		}

		lineIsBullet = true

		return
	}

	if m := logLineContinuation.FindStringSubmatch(line); m != nil {
		if status.previousLineWasBullet {
			lineIsBullet = true
		} else {
			status.fail(lineNumber, "log line continuation has not bullet to continue")
		}

		return
	}

	if blankLine.MatchString(line) {
		if status.previousLineWasBlank {
			status.fail(lineNumber, "consecutive blank lines")
		}

		lineIsBlank = true

		return
	}

	if whitespaceLine.MatchString(line) {
		status.fail(lineNumber, "line looks blank, but has spaces in it")

		return
	}

	status.fail(lineNumber, "badly formatted line: it must be blank, start with a space or bullet (\" * \"), or a heading of the form \"version  date\"")
}
