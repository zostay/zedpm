package goals

import (
	"context"
	"time"

	"github.com/zostay/zedpm/plugin"
	"github.com/zostay/zedpm/storage"
)

const (
	PropertyExportPrefix = storage.ExportPrefix

	PropertyReleaseDescription = "release.description"
	PropertyReleaseVersion     = "release.version"
	PropertyReleaseDate        = "release.date"

	PropertyLintPreRelease = "lint.prerelease"
	PropertyLintRelease    = "lint.release"

	PropertyInfoVersion = "info.version"
)

func SetPropertyReleaseDescription(
	ctx context.Context,
	description string,
) {
	plugin.Set(ctx, PropertyReleaseDescription, description)
}

func GetPropertyInfoVersion(ctx context.Context) string {
	return plugin.GetString(ctx, PropertyInfoVersion)
}

func GetPropertyReleaseDescription(ctx context.Context) string {
	return plugin.GetString(ctx, PropertyReleaseDescription)
}

func SetPropertyReleaseVersion(ctx context.Context, version string) {
	plugin.Set(ctx, PropertyReleaseVersion, version)
}

func GetPropertyReleaseVersion(ctx context.Context) string {
	return plugin.GetString(ctx, PropertyReleaseVersion)
}

func SetPropertyReleaseDate(ctx context.Context, date time.Time) {
	plugin.Set(ctx, PropertyReleaseDate, date)
}

func GetPropertyReleaseDate(ctx context.Context) time.Time {
	return plugin.GetTime(ctx, PropertyReleaseDate)
}

func ExportPropertyName(
	ctx context.Context,
	propertyName string,
) {
	plugin.Set(ctx, PropertyExportPrefix+propertyName, true)
}

func SetPropertyLintPreRelease(ctx context.Context, value bool) {
	plugin.Set(ctx, PropertyLintPreRelease, value)
}

func GetPropertyLintPreRelease(ctx context.Context) bool {
	return plugin.GetBool(ctx, PropertyLintPreRelease)
}

func SetPropertyLintRelease(ctx context.Context, value bool) {
	plugin.Set(ctx, PropertyLintRelease, value)
}

func GetPropertyLintRelease(ctx context.Context) bool {
	return plugin.GetBool(ctx, PropertyLintRelease)
}
