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

	if version := goals.GetPropertyReleaseVersion(ctx); version != "" {
		return defaultReleaseNamePrefix + version, nil
	}

	return "", fmt.Errorf("missing required properties %q or %q", PropertyGithubReleaseName, goals.PropertyReleaseVersion)
}
