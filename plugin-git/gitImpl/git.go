package gitImpl

import (
	"context"

	"github.com/zostay/zedpm/plugin"
)

const (
	PropertyGitIgnoreDirty = "git.ignoreDirty"
)

// GetPropertyGitIgnoreDirty gets the setting, that when true, allows the minting
// to proceed even when the directory is dirty.
func GetPropertyGitIgnoreDirty(ctx context.Context) bool {
	return plugin.GetBool(ctx, PropertyGitIgnoreDirty)
}
