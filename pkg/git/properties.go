package git

import (
	"context"

	"github.com/zostay/zedpm/format"
	"github.com/zostay/zedpm/pkg/goals"
	"github.com/zostay/zedpm/plugin"
)

const (
	PropertyGitReleaseTag    = "git.release.tag"
	PropertyGitReleaseBranch = "git.release.branch"
)

const (
	defaultReleaseBranchPrefix = "release-v"
	defaultReleaseTagPrefix    = "v"
)

func GetPropertyGitReleaseTag(ctx context.Context) (string, error) {
	if plugin.IsSet(ctx, PropertyGitReleaseTag) {
		return plugin.GetString(ctx, PropertyGitReleaseTag), nil
	}

	version, err := goals.GetPropertyReleaseVersion(ctx)
	if err != nil {
		return "", format.WrapErr(err, "unable to find or create a value for %q", PropertyGitReleaseTag)
	}

	tagName := defaultReleaseTagPrefix + version
	return tagName, nil
}

func GetPropertyGitReleaseBranch(ctx context.Context) (string, error) {
	if plugin.IsSet(ctx, PropertyGitReleaseBranch) {
		return plugin.GetString(ctx, PropertyGitReleaseBranch), nil
	}

	version, err := goals.GetPropertyReleaseVersion(ctx)
	if err != nil {
		return "", format.WrapErr(err, "unable to find or create a value for %q", PropertyGitReleaseBranch)
	}

	branchName := defaultReleaseBranchPrefix + version
	return branchName, nil
}
