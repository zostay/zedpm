package changelogImpl

import (
	"context"

	"github.com/zostay/zedpm/plugin"
)

// LintChangelogTask implements the /lint/changelog task.
type LintChangelogTask struct {
	plugin.TaskBoilerplate
}

// Run prepares the system to run the LintChangelog operation.
func (t *LintChangelogTask) Run(_ context.Context) (plugin.Operations, error) {
	return plugin.Operations{
		{
			Order:  50,
			Action: plugin.OperationFunc(LintChangelog),
		},
	}, nil
}
