package goals

import (
	"context"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"

	"github.com/zostay/zedpm/storage"
)

// OutputFormatter is a function that can output a storage.KV to the given io.Writer.
type OutputFormatter func(io.Writer, storage.KV) error

// OutputFormats defines the available output formats.
var OutputFormats = map[string]OutputFormatter{
	"properties": WriteOutProperties,
	"yaml":       WriteOutYaml,
}

// DefaultInfoOutputFormatter is the default output format.
var DefaultInfoOutputFormatter = WriteOutProperties

// InfoOutputFormatter tries to determine which output formatter to use based
// upon the info.outputFormat property. If the property is not set or set to an
// invalid value, this will return DefaultInfoOutputFormatter.
func InfoOutputFormatter(ctx context.Context) OutputFormatter {
	format := GetPropertyInfoOutputFormat(ctx)
	formatter := OutputFormats[format]
	if formatter != nil {
		return formatter
	}
	return DefaultInfoOutputFormatter
}

// WriteOutProperties outputs the given storage values as a properties list.
func WriteOutProperties(w io.Writer, values storage.KV) error {
	for _, key := range values.AllKeys() {
		_, err := fmt.Fprintf(w, "%s = %#v\n", key, values.Get(key))
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteOutYaml outputs the given storage values in YAML format.
func WriteOutYaml(w io.Writer, values storage.KV) error {
	enc := yaml.NewEncoder(w)
	return enc.Encode(values.AllSettings())
}
