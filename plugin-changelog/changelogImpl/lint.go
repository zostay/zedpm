package changelogImpl

import (
	"context"

	"github.com/zostay/zedpm/plugin"
)

type LintChangelogTask struct {
	plugin.TaskBoilerplate
}

func (t *LintChangelogTask) Run(ctx context.Context) (plugin.Operations, error) {
	return plugin.Operations{
		{
			Order:  50,
			Action: plugin.OperationFunc(LintChangelog),
		},
	}, nil
}
