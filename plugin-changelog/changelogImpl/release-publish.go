package changelogImpl

import (
	"context"
	"io"

	"github.com/zostay/zedpm/format"
	"github.com/zostay/zedpm/pkg/changes"
	"github.com/zostay/zedpm/pkg/goals"
	"github.com/zostay/zedpm/plugin"
)

// ReleasePublishTask implements the /release/publish/changelog task.
type ReleasePublishTask struct {
	plugin.TaskBoilerplate
}

// CaptureChangesInfo loads the bullets for the changelog section relevant to
// this release into the process configuration for use when creating the release
// later.
func (f *ReleasePublishTask) CaptureChangesInfo(ctx context.Context) error {
	version, err := goals.GetPropertyReleaseVersion(ctx)
	if err != nil {
		return err
	}

	vstring := "v" + version
	changelog := GetPropertyChangelogFile(ctx)
	cr, err := changes.ExtractSection(changelog, vstring)
	if err != nil {
		return format.WrapErr(err, "unable to get log of changes")
	}

	chgs, err := io.ReadAll(cr)
	if err != nil {
		return format.WrapErr(err, "unable to read log of changes")
	}

	plugin.Set(ctx, goals.PropertyReleaseDescription, string(chgs))

	plugin.Logger(ctx,
		"version", version,
	).Info("Captured release description from changelog for version %q.", version)

	return nil
}

// Check executes CaptureChangesInfo to get the latest changes and save them for
// release.description for use by other plugins to finish the release process.
func (f *ReleasePublishTask) Check(ctx context.Context) error {
	return f.CaptureChangesInfo(ctx)
}
