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

	if version := goals.GetPropertyReleaseVersion(ctx); version != "" {
		tagName := defaultReleaseTagPrefix + version
		return tagName, nil
	}

	return "", fmt.Errorf("missing required %q or %q settings", PropertyGitReleaseTag, goals.PropertyReleaseVersion)
}

func GetPropertyGitReleaseBranch(ctx context.Context) (string, error) {
	if plugin.IsSet(ctx, PropertyGitReleaseBranch) {
		return plugin.GetString(ctx, PropertyGitReleaseBranch), nil
	}

	if version := goals.GetPropertyReleaseVersion(ctx); version != "" {
		branchName := defaultReleaseBranchPrefix + version
		return branchName, nil
	}

	return "", fmt.Errorf("missing required %q or %q settings", PropertyGitReleaseBranch, goals.PropertyReleaseVersion)
}
