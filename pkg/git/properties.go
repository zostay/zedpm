package git

import (
	"context"
	"fmt"

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
		return "", fmt.Errorf("unable to find or create a value for %q: %w", PropertyGitReleaseTag, err)
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
		return "", fmt.Errorf("unable to find or create a value for %q: %w", PropertyGitReleaseBranch, err)
	}

	branchName := defaultReleaseBranchPrefix + version
	return branchName, nil
}
