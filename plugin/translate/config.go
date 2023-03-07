package translate

import (
	"github.com/zostay/zedpm/plugin/api"
	"github.com/zostay/zedpm/storage"
)

// APIConfigToKV translates an api.Config object into a storage.KVMem object.
func APIConfigToKV(in *api.Config) *storage.KVMem {
	out := storage.New()

	for k, v := range in.GetValues() {
		out.Set(k, v)
	}

	return out
}

// KVToAPIConfig translates a storage.KV object into an api.Config object.
func KVToAPIConfig(in storage.KV) *api.Config {
	return &api.Config{Values: KVToStringMapString(in)}
}

// KVToStringMapString translates a storage.KV object into a map[string]string.
func KVToStringMapString(in storage.KV) map[string]string {
	keys := in.AllKeys()
	out := make(map[string]string, len(keys))

	for _, k := range keys {
		out[k] = in.GetString(k)
	}
	return out
}
