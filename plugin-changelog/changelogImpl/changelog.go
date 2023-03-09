package changelogImpl

import (
	"context"
	"fmt"
	"os"

	"github.com/zostay/zedpm/format"
	"github.com/zostay/zedpm/pkg/changes"
	"github.com/zostay/zedpm/pkg/goals"
	"github.com/zostay/zedpm/plugin"
)

const (
	PropertyChangelogFile = "changlog.file"
)

// DefaultChangelog is the changelog file path to use when none is configured.
const DefaultChangelog = "Changes.md"

// GetPropertyChangelogFile gets the name of the changelog file from the configuration or
// returns the default value.
func GetPropertyChangelogFile(ctx context.Context) string {
	if plugin.IsSet(ctx, PropertyChangelogFile) {
		return plugin.GetString(ctx, PropertyChangelogFile)
	}
	return DefaultChangelog
}

// CheckMode examines the context to determine what mode to use when performing
// linter checks.
//
// If lint.release is set, then it uses changes.CheckRelease. If lint.prerelease
// is set, then it uses changes.CheckPreRelease. Otherwise, it uses
// changes.CheckStandard.
func CheckMode(ctx context.Context) changes.CheckMode {
	var mode changes.CheckMode
	switch {
	case goals.GetPropertyLintRelease(ctx):
		mode = changes.CheckRelease
	case goals.GetPropertyLintPreRelease(ctx):
		mode = changes.CheckPreRelease
	default:
		mode = changes.CheckStandard
	}
	return mode
}

// LintChangelog performs a check to ensure the changelog is ready for release.
func LintChangelog(ctx context.Context) error {
	changelog, err := os.Open(GetPropertyChangelogFile(ctx))
	if err != nil {
		return format.WrapErr(err, "unable to open Changes file")
	}

	linter := changes.NewLinter(changelog, CheckMode(ctx))
	err = linter.Check()
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
