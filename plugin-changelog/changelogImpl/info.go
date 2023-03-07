package changelogImpl

import (
	"context"
	"fmt"
	"io"

	"github.com/zostay/zedpm/changes"
	"github.com/zostay/zedpm/plugin"
	"github.com/zostay/zedpm/plugin-goals/pkg/goals"
)

// InfoChangelogTask implements the /info/release/description task, which describes either
// the current release or the release named in info.version.
type InfoChangelogTask struct {
	plugin.TaskBoilerplate
}

// ExtractChangelog does the work of extracting a changelog section for a single
// version. If no version is specified in info.version, the first (latest)
// version is used.
func (t *InfoChangelogTask) ExtractChangelog(ctx context.Context) error {
	version := goals.GetPropertyInfoVersion(ctx)
	r, err := changes.ExtractSection(GetPropertyChangelogFile(ctx), version)
	if err != nil {
		return fmt.Errorf("failed to read changelog section: %w", err)
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read changelog data: %w", err)
	}

	goals.SetPropertyReleaseDescription(ctx, string(data))
	goals.ExportPropertyName(ctx, goals.PropertyReleaseDescription)

	return nil
}

// Run prepares the ExtractChangelog operation to run.
func (t *InfoChangelogTask) Run(context.Context) (plugin.Operations, error) {
	return plugin.Operations{
		{
			Order:  50,
			Action: plugin.OperationFunc(t.ExtractChangelog),
		},
	}, nil
}
