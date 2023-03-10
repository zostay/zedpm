package git

import (
	"context"

	"github.com/zostay/zedpm/format"
	"github.com/zostay/zedpm/pkg/goals"
	"github.com/zostay/zedpm/plugin"
)

const (
	PropertyReleaseBranchPrefix = "release-v"
	DefaultGitTargetBranch      = "master"

	PropertyGitReleaseTag    = "git.release.tag"
	PropertyGitReleaseBranch = "git.release.branch"
	PropertyGitTargetBranch  = "git.target.branch"
)

func GetPropertyGitReleaseTag(ctx context.Context) (string, error) {
	tag := plugin.GetString(ctx, PropertyGitReleaseTag)
	if tag != "" {
		return tag, nil
	}

	return goals.GetPropertyReleaseTag(ctx)
}

func GetPropertyGitReleaseBranch(ctx context.Context) (string, error) {
	branch := plugin.GetString(ctx, PropertyGitReleaseBranch)
	if branch != "" {
		return branch, nil
	}

	version, err := goals.GetPropertyReleaseVersion(ctx)
	if err != nil {
		return "", format.WrapErr(err, "unable to find or compute %q setting", PropertyGitReleaseBranch)
	}

	return PropertyReleaseBranchPrefix + version, nil
}

func GetPropertyGitTargetBranch(ctx context.Context) string {
	branch := plugin.GetString(ctx, PropertyGitTargetBranch)
	if branch != "" {
		return branch
	}

	return DefaultGitTargetBranch
}
