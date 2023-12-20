package changelogImpl

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/zostay/zedpm/format"
	"github.com/zostay/zedpm/pkg/goals"
	"github.com/zostay/zedpm/plugin"
)

// ReleaseMintTask is the implementation of the /release/mint/changelog task.
type ReleaseMintTask struct {
	plugin.TaskBoilerplate
}

// FixupChangelog alters the changelog to prepare it for release.
func (s *ReleaseMintTask) FixupChangelog(ctx context.Context) error {
	r, err := os.Open(GetPropertyChangelogFile(ctx))
	if err != nil {
		return format.WrapErr(err, "unable to open %s", GetPropertyChangelogFile(ctx))
	}

	newChangelog := GetPropertyChangelogFile(ctx) + ".new"

	w, err := os.Create(newChangelog)
	if err != nil {
		return format.WrapErr(err, "unable to create %s", newChangelog)
	}

	plugin.ForCleanup(ctx, func() { _ = os.Remove(newChangelog) })

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		if line == "WIP" || line == "WIP  TBD" {
			version, err := goals.GetPropertyReleaseVersion(ctx)
			if err != nil {
				return err
			}

			todayTime := goals.GetPropertyReleaseDate(ctx)

			today := todayTime.Format("2006-01-02")
			_, _ = fmt.Fprintf(w, "v%s  %s\n", version, today)
		} else {
			_, _ = fmt.Fprintln(w, line)
		}
	}

	_ = r.Close()
	err = w.Close()
	if err != nil {
		return format.WrapErr(err, "unable to close %s", newChangelog)
	}

	err = os.Rename(newChangelog, GetPropertyChangelogFile(ctx))
	if err != nil {
		return format.WrapErr(err, "unable to overwrite %s with %s", GetPropertyChangelogFile(ctx), newChangelog)
	}

	plugin.ToAdd(ctx, GetPropertyChangelogFile(ctx))

	plugin.Logger(ctx,
		"changelog", GetPropertyChangelogFile(ctx),
	).Info("Applied changes to changelog to fixup for release.")

	return nil
}

// Check lints the changelog for pre-release.
func (s *ReleaseMintTask) Check(ctx context.Context) error {
	goals.SetPropertyLintPreRelease(ctx, true)
	goals.SetPropertyLintRelease(ctx, false)

	return LintChangelog(ctx)
}

// LintChangelogRelease lints the changelog for release.
func (s *ReleaseMintTask) LintChangelogRelease(ctx context.Context) error {
	goals.SetPropertyLintPreRelease(ctx, false)
	goals.SetPropertyLintRelease(ctx, true)

	return LintChangelog(ctx)
}

// Run prepares the FixupChangelog and the final LintChangelogRelease
// operations.
func (s *ReleaseMintTask) Run(context.Context) (plugin.Operations, error) {
	return plugin.Operations{
		{
			Order:  50,
			Action: plugin.OperationFunc(s.FixupChangelog),
		},
		{
			Order:  55,
			Action: plugin.OperationFunc(s.LintChangelogRelease),
		},
	}, nil
}
