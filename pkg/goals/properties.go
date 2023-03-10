package goals

import (
	"context"
	"fmt"
	"time"

	"github.com/zostay/zedpm/format"
	"github.com/zostay/zedpm/pkg/storage"
	"github.com/zostay/zedpm/plugin"
)

const (
	PropertyReleaseTagPrefix = "v"
	PropertyExportPrefix     = storage.ExportPrefix

	PropertyReleaseDescription = "release.description"
	PropertyReleaseVersion     = "release.version"
	PropertyReleaseDate        = "release.date"
	PropertyReleaseTag         = "release.tag"

	PropertyLintPreRelease = "lint.prerelease"
	PropertyLintRelease    = "lint.release"

	PropertyInfoVersion      = "info.version"
	PropertyInfoOutputFormat = "info.outputFormat"
	PropertyInfoOutputAll    = "info.outputAll"

	DefaultInfoOutputFormat = "properties"
)

// SetPropertyReleaseDescription sets the release.description property to the
// given description.
func SetPropertyReleaseDescription(
	ctx context.Context,
	description string,
) {
	plugin.Set(ctx, PropertyReleaseDescription, description)
}

// GetPropertyInfoVersion returns the value of info.version.
func GetPropertyInfoVersion(ctx context.Context) string {
	return plugin.GetString(ctx, PropertyInfoVersion)
}

// GetPropertyInfoOutputFormat returns the value of info.outputFormat. If not
// found, this will return DefaultInfoOutputFormat.
func GetPropertyInfoOutputFormat(ctx context.Context) string {
	if format := plugin.GetString(ctx, PropertyInfoOutputFormat); format != "" {
		return format
	}
	return DefaultInfoOutputFormat
}

// GetPropertyInfoOutputAll returns the value of info.outputAll.
func GetPropertyInfoOutputAll(ctx context.Context) bool {
	return plugin.GetBool(ctx, PropertyInfoOutputAll)
}

// GetPropertyReleaseDescription returns the value of release.description.
func GetPropertyReleaseDescription(ctx context.Context) string {
	desc := plugin.GetString(ctx, PropertyReleaseDescription)
	if desc == "" {
		desc = "No description provided."
	}
	return desc
}

// SetPropertyReleaseVersion sets the value of release.version.
func SetPropertyReleaseVersion(ctx context.Context, version string) {
	plugin.Set(ctx, PropertyReleaseVersion, version)
}

// GetPropertyReleaseVersion gets the value of release.version.
func GetPropertyReleaseVersion(ctx context.Context) (string, error) {
	version := plugin.GetString(ctx, PropertyReleaseVersion)
	if version != "" {
		return version, nil
	}

	return "", fmt.Errorf("%q is not defined", PropertyReleaseVersion)
}

// SetPropertyReleaseDate sets the value of release.date.
func SetPropertyReleaseDate(ctx context.Context, date time.Time) {
	plugin.Set(ctx, PropertyReleaseDate, date)
}

// GetPropertyReleaseDate gets the value of release.date.
func GetPropertyReleaseDate(ctx context.Context) time.Time {
	if date := plugin.GetTime(ctx, PropertyReleaseDate); !date.IsZero() {
		return date
	}

	return time.Now()
}

// ExportPropertyName sets the given property name with teh PropertyExprotPrefix
// so that it will be rendered by the /info/display task when the info goal is
// complete.
func ExportPropertyName(
	ctx context.Context,
	propertyName string,
) {
	plugin.Set(ctx, PropertyExportPrefix+propertyName, true)
}

// SetPropertyLintPreRelease sets the lint.prerelease property.
func SetPropertyLintPreRelease(ctx context.Context, value bool) {
	plugin.Set(ctx, PropertyLintPreRelease, value)
}

// GetPropertyLintPreRelease gets the lint.prerelease property.
func GetPropertyLintPreRelease(ctx context.Context) bool {
	return plugin.GetBool(ctx, PropertyLintPreRelease)
}

// SetPropertyLintRelease sets the lint.release property.
func SetPropertyLintRelease(ctx context.Context, value bool) {
	plugin.Set(ctx, PropertyLintRelease, value)
}

// GetPropertyLintRelease gets the lint.release property.
func GetPropertyLintRelease(ctx context.Context) bool {
	return plugin.GetBool(ctx, PropertyLintRelease)
}

// GetPropertyReleaseTag gets the release.tag property.
func GetPropertyReleaseTag(ctx context.Context) (string, error) {
	if plugin.IsSet(ctx, PropertyReleaseTag) {
		return plugin.GetString(ctx, PropertyReleaseTag), nil
	}

	version, err := GetPropertyReleaseVersion(ctx)
	if err != nil {
		return "", format.WrapErr(err, "unable to get or create a value for %q", PropertyReleaseTag)
	}

	return PropertyReleaseTagPrefix + version, nil
}
