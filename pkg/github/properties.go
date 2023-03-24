package github

import (
	"context"

	"github.com/zostay/zedpm/format"
	"github.com/zostay/zedpm/pkg/goals"
	"github.com/zostay/zedpm/plugin"
)

const (
	PropertyGithubReleaseName = "github.release.name"
	PropertyGithubOwner       = "github.owner"
	PropertyGithubProject     = "github.project"
)

const defaultReleaseNamePrefix = "Release v"

func GetPropertyGithubReleaseName(ctx context.Context) (string, error) {
	if plugin.IsSet(ctx, PropertyGithubReleaseName) {
		return plugin.GetString(ctx, PropertyGithubReleaseName), nil
	}

	version, err := goals.GetPropertyReleaseVersion(ctx)
	if err != nil {
		return "", format.WrapErr(err, "unable to get or create a value for %q", PropertyGithubReleaseName)
	}

	return defaultReleaseNamePrefix + version, nil
}

func GetPropertyGithubOwner(ctx context.Context) string {
	return plugin.GetString(ctx, PropertyGithubOwner)
}

func GetPropertyGithubProject(ctx context.Context) string {
	return plugin.GetString(ctx, PropertyGithubProject)
}
