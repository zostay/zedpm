package goalsImpl

import (
	"context"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"

	"github.com/zostay/zedpm/plugin"
	"github.com/zostay/zedpm/storage"
)

type OutputFormatter func(io.Writer, storage.KV) error

const InfoOutputFormatKey = "info.outputFormat"
const InfoOutputAllKey = "info.outputAll"

var outputFormats = map[string]OutputFormatter{
	"properties": WriteOutProperties,
	"yaml":       WriteOutYaml,
}

var DefaultInfoOutputFormatter = WriteOutProperties

func InfoOutputFormatter(ctx context.Context) OutputFormatter {
	format := plugin.GetString(ctx, InfoOutputFormatKey)
	formatter := outputFormats[format]
	if formatter != nil {
		return formatter
	}
	return DefaultInfoOutputFormatter
}

func InfoOutputAll(ctx context.Context) bool {
	return plugin.GetBool(ctx, InfoOutputAllKey)
}

func WriteOutProperties(w io.Writer, values storage.KV) error {
	for _, key := range values.AllKeys() {
		// TODO is key.subkey.subsubkey....=value the best output format?
		_, err := fmt.Fprintf(w, "%s = %#v\n", key, values.Get(key))
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteOutYaml(w io.Writer, values storage.KV) error {
	enc := yaml.NewEncoder(w)
	return enc.Encode(values.AllSettings())
}
