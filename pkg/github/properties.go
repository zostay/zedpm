package github

import (
	"context"
	"fmt"

	"github.com/zostay/zedpm/pkg/goals"
	"github.com/zostay/zedpm/plugin"
)

const PropertyGithubReleaseName = "github.release.name"

const defaultReleaseNamePrefix = "Release v"

func GetPropertyGithubReleaseName(ctx context.Context) (string, error) {
	if plugin.IsSet(ctx, PropertyGithubReleaseName) {
		return plugin.GetString(ctx, PropertyGithubReleaseName), nil
	}

	version, err := goals.GetPropertyReleaseVersion(ctx)
	if err != nil {
		return "", fmt.Errorf("unable to get or create a value for %q: %w", PropertyGithubReleaseName, err)
	}

	return defaultReleaseNamePrefix + version, nil
}
