package changelogImpl

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/zostay/zedpm/plugin"
	"github.com/zostay/zedpm/plugin-goals/pkg/goals"
)

// ReleaseMintTask is the implementation of the /release/mint/changelog task.
type ReleaseMintTask struct {
	plugin.TaskBoilerplate
}

// FixupChangelog alters the changelog to prepare it for release.
func (s *ReleaseMintTask) FixupChangelog(ctx context.Context) error {
	r, err := os.Open(GetPropertyChangelogFile(ctx))
	if err != nil {
		return fmt.Errorf("unable to open %s: %w", GetPropertyChangelogFile(ctx), err)
	}

	newChangelog := GetPropertyChangelogFile(ctx) + ".new"

	w, err := os.Create(newChangelog)
	if err != nil {
		return fmt.Errorf("unable to create %s: %w", newChangelog, err)
	}

	plugin.ForCleanup(ctx, func() { _ = os.Remove(newChangelog) })

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		if line == "WIP" || line == "WIP  TBD" {
			version := goals.GetPropertyReleaseVersion(ctx)
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
		return fmt.Errorf("unable to close %s: %w", newChangelog, err)
	}

	err = os.Rename(newChangelog, GetPropertyChangelogFile(ctx))
	if err != nil {
		return fmt.Errorf("unable to overwrite %s with %s: %w", GetPropertyChangelogFile(ctx), newChangelog, err)
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
