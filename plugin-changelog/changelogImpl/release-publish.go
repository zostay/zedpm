package changelogImpl

import (
	"context"
	"fmt"
	"io"

	"github.com/zostay/zedpm/changes"
	"github.com/zostay/zedpm/plugin"
	"github.com/zostay/zedpm/plugin-goals/pkg/goals"
)

type ReleasePublishTask struct {
	plugin.TaskBoilerplate

	Changelog string
}

// CaptureChangesInfo loads the bullets for the changelog section relevant to
// this release into the process configuration for use when creating the release
// later.
func (f *ReleasePublishTask) CaptureChangesInfo(ctx context.Context) error {
	version := plugin.GetString(ctx, "release.version")
	vstring := "v" + version
	cr, err := changes.ExtractSection(f.Changelog, vstring)
	if err != nil {
		return fmt.Errorf("unable to get log of changes: %w", err)
	}

	chgs, err := io.ReadAll(cr)
	if err != nil {
		return fmt.Errorf("unable to read log of changes: %w", err)
	}

	plugin.Set(ctx, goals.PropertyReleaseDescription, string(chgs))

	plugin.Logger(ctx,
		"version", version,
	).Info("Captured release description from changelog for version %q.", version)

	return nil
}

func (f *ReleasePublishTask) Check(ctx context.Context) error {
	if f.Changelog == "" {
		return fmt.Errorf("missing changelog location in paths config")
	}

	return f.CaptureChangesInfo(ctx)
}
