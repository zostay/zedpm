package storage

import "time"

// KV describes a key-value store interface that is used to hold configuration
// properties used by the various tools in zedpm. The internal storage of a KV
// is a hierarchical map of maps. This way deeply nested structures can be held.
// The deeply nested values can be retrieved using keys in dot-format.
//
// For example, setting a value with a key named "git.release.tag" will ensure
// that the top-level map contains a key named "git" that points to
// map[string]any pointer. That map object will contain a key named "release"
// that points to another map[string]any pointer. This inner map will then have
// a key named "tag" which points to the value set.
//
// The KV allows any kind of value to be set, but all values should coerce to
// string in order to allow them to be communicated over the wire between
// plugins.
type KV interface {
	// AllKeys returns a string list representing all the key names stored in
	// the KV. Keys in this list will be in dot-format.
	AllKeys() []string

	// AllSettings returns a hierarchical map of maps containing all the
	// properties that are set on this KV.
	AllSettings() map[string]any

	// AllSettingsStrings returns a flat map of key/value pairs containing all
	// the properties that are set on this KV. Keys in this map are in
	// dot-separated format.
	AllSettingsStrings() map[string]string

	// Get a value from the KV.
	Get(string) any

	// GetBool gets a boolean value from the KV.
	GetBool(string) bool

	// GetDuration gets a time.Duration value from the KV.
	GetDuration(string) time.Duration

	// GetFloat64 gets a float64 value from the KV.
	GetFloat64(string) float64

	// GetInt gets an int value from the KV.
	GetInt(string) int

	// GetInt32 gets an int32 value from the KV.
	GetInt32(string) int32

	// GetInt64 gets an int64 value from the KV.
	GetInt64(string) int64

	// GetIntSlice gets a slice of ints value from the KV.
	GetIntSlice(string) []int

	// GetString gets a string value from the KV.
	GetString(string) string

	// GetStringMap gets a map of string/any pairs from the KV.
	GetStringMap(string) map[string]any

	// GetStringMapString gets a map of string/string pairs from the KV.
	GetStringMapString(string) map[string]string

	// GetStringMapStringSlice gets a map of string/slice of strings pairs from
	// the KV.
	GetStringMapStringSlice(string) map[string][]string

	// GetStringSlice gets a slice of strings from the KV.
	GetStringSlice(string) []string

	// GetTime gets a time.Time value from the KV.
	GetTime(string) time.Time

	// GetUint gets a uint value from the KV.
	GetUint(string) uint

	// GetUint16 gets a uint16 value from the KV.
	GetUint16(string) uint16

	// GetUint32 gets a uint16 value from the KV.
	GetUint32(string) uint32

	// GetUint64 gets a uint16 value from the KV.
	GetUint64(string) uint64

	// Sub returns a KV that access keys under a subtree.
	Sub(string) KV

	// IsSet returns true if the KV contains a value for the given key.
	IsSet(string) bool

	// Clear deletes all values from the KV.
	Clear()

	// Set sets a value with the given key on the KV.
	Set(string, any)

	// Update sets all the values in the given map on the KV. This overwrites
	// existing keys. It does not perform a merge.
	Update(map[string]any)

	// UpdateStrings sets all the values in the given map on the KV. This will
	// perform a merge operation rather than overwrite.
	UpdateStrings(map[string]string)

	// RegisterAlias establishes a relationship between an old key name and a
	// new key name.
	RegisterAlias(string, string)
}

// TODO I created Requirements a bit too early. Does this make sense or shall it be deleted as deadcode?

// Requirements is an additional layer that can be added to a KV to apply
// requirements to KVs.
type Requirements interface {
	// MarkRequired marks the name key as a required value that must be set.
	MarkRequired(string)

	// IsRequired returns true if the given key has been marked as required.
	IsRequired(string) bool

	// IsPassingRequirements returns true if all the required keys are set.
	IsPassingRequirements() bool

	// MissingRequirements returns the names of all the keys that are required,
	// but not currently not set.
	MissingRequirements() []string
}
