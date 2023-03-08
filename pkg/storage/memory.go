package storage

import (
	"container/list"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cast"
)

// Verifies that KVMem is a KV and implements Requirements.
var (
	_ KV           = &KVMem{}
	_ Requirements = &KVMem{}
)

// KVMem is the base building block of the storage package. It provides a basic
// in-memory store of a hierarchy of key/value pairs.
type KVMem struct {
	prefix       string
	values       map[string]any
	requirements map[string]struct{}
	aliases      map[string]string
}

// New returns a new, empty KVMem.
func New() *KVMem {
	return &KVMem{
		values:       make(map[string]any, 10),
		requirements: make(map[string]struct{}, 10),
	}
}

// RO returns a KVCfg, which is a read-only version of the KVMem.
func (m *KVMem) RO() *KVCfg {
	return &KVCfg{*m}
}

// prefixKey adds the prefix to the value in case this KVMem was constructed via
// Sub.
func (m *KVMem) prefixKey(key string) string {
	if m.prefix != "" {
		return m.prefix + "." + key
	}
	return key
}

// resolveAlias prefixes the given key and changes the name of the key to the
// target key name.
func (m *KVMem) resolveAlias(key string) string {
	prefixedKey := m.prefixKey(key)
	if otherKey, isAliased := m.aliases[prefixedKey]; isAliased {
		return m.prefixKey(otherKey)
	}
	return prefixedKey
}

// key transforms the key to lower case, adds the prefix, and then resolves an
// alias, if any.
func (m *KVMem) key(key string) string {
	key = strings.ToLower(key)
	return m.resolveAlias(key)
}

// splitKey splits the key up into its components by splitting by ".".
func (m *KVMem) splitKey(key string) []string {
	return strings.Split(m.key(key), ".")
}

// getExists does the work of searching down the hierarchy with the given key
// for its value. If the value can't be found, it returns nil and false. If a
// value can be found, it returns the value (which may be nil or anything else)
// and true.
func (m *KVMem) getExists(key string) (any, bool) {
	keys := m.splitKey(key)
	n := m.values
	for i, k := range keys {
		if n == nil {
			return n, false
		}

		val, exists := n[k]
		if !exists {
			return nil, false
		} else if i == len(keys)-1 {
			// last key being resolved, return the value
			return val, exists
		}

		// not the last key, so we need map[string]any
		switch v := val.(type) {
		case map[string]any:
			n = v
		default:
			// we can't go further, but they want us to, so the key they
			// want doesn't exist
			return nil, false
		}
	}
	return n, true
}

// get retrieves the value from the given key or nil if it does not exist.
func (m *KVMem) get(key string) any {
	v, _ := m.getExists(key)
	return v
}

// set searches down through the tree, setting any parent keys to empty
// map[string]any and then setting the final key to the given value.
func (m *KVMem) set(key string, value any) {
	keys := m.splitKey(key)
	var lastKey string
	keys, lastKey = keys[:len(keys)-1], keys[len(keys)-1]
	values := m.values
	for _, k := range keys {
		xVals := values[k]
		if vals, isStringMap := xVals.(map[string]any); isStringMap {
			values = vals
		} else {
			vals := map[string]any{}
			values, values[k] = vals, vals
		}
	}
	values[lastKey] = value
}

// keys is a generic tool for grabbing the keys off of a map.
func keys[T any](v map[string]T) []string {
	ks := make([]string, 0, len(v))
	for k := range v {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

// AllKeys returns a slice of strings that point to all the end values. It does
// not include any key that would point to a map of more values somewhere down
// the tree.
func (m *KVMem) AllKeys() []string {
	type openItem struct {
		prefix string
		keys   []string
		in     map[string]any
	}

	out := make([]string, 0, len(m.values))
	openList := list.New()
	openList.PushFront(&openItem{
		prefix: "",
		keys:   keys[any](m.values),
		in:     m.values,
	})
	for openList.Len() > 0 {
		el := openList.Front()
		item := el.Value.(*openItem)

		if len(item.keys) == 0 {
			openList.Remove(el)
			continue
		}

		var key string
		key, item.keys = item.keys[0], item.keys[1:]
		value := item.in[key]

		if item.prefix != "" {
			key = item.prefix + "." + key
		}

		if nextIn, isStringMap := value.(map[string]any); isStringMap {
			openList.PushFront(&openItem{
				prefix: key,
				keys:   keys[any](nextIn),
				in:     nextIn,
			})
			continue
		}

		// only append the key name if it's a value, not a nested map[string]any
		out = append(out, key)
	}

	out = append(out, keys[string](m.aliases)...)
	sort.Strings(out)

	return out
}

// AllSettings returns a deep copy of the internal storage of KVMem.
func (m *KVMem) AllSettings() map[string]any {
	out := make(map[string]any, len(m.values))
	for k, v := range m.values {
		out[k] = v
	}

	tmp := &KVMem{values: out}
	for alias, key := range m.aliases {
		tmp.set(alias, tmp.get(key))
	}

	return tmp.values
}

// AllSettingsStrings returns a flattened copy of the internal storage of KVMem.
func (m *KVMem) AllSettingsStrings() map[string]string {
	keys := m.AllKeys()
	out := make(map[string]string, len(keys))
	for _, k := range keys {
		out[k] = m.GetString(k)
	}
	return out
}

// Get returns the value for the given key or nil.
func (m *KVMem) Get(key string) any {
	return m.get(key)
}

// GetBool returns the value of the given key as a boolean or false.
func (m *KVMem) GetBool(key string) bool {
	return cast.ToBool(m.get(key))
}

// GetDuration returns ehe value of the given key as a duration or zero.
func (m *KVMem) GetDuration(key string) time.Duration {
	return cast.ToDuration(m.get(key))
}

// GetFloat64 returns the value of the given key as a floating point value or
// zero.
func (m *KVMem) GetFloat64(key string) float64 {
	return cast.ToFloat64(m.get(key))
}

// GetInt returns the value of the given key as an integer or zero.
func (m *KVMem) GetInt(key string) int {
	return cast.ToInt(m.get(key))
}

// GetInt32 returns the value of the given key as a 32-bit integer or zero.
func (m *KVMem) GetInt32(key string) int32 {
	return cast.ToInt32(m.get(key))
}

// GetInt64 returns the value of the given key as a 64-bit integer or zero.
func (m *KVMem) GetInt64(key string) int64 {
	return cast.ToInt64(m.get(key))
}

// GetIntSlice returns the value of the given key as an integer slice or nil.
func (m *KVMem) GetIntSlice(key string) []int {
	return cast.ToIntSlice(m.get(key))
}

// GetString returns the value of the given key as a string or an empty string.
func (m *KVMem) GetString(key string) string {
	return cast.ToString(m.get(key))
}

// GetStringMap returns the value of the given key as a map of strings to any or
// nil.
func (m *KVMem) GetStringMap(key string) map[string]any {
	return cast.ToStringMap(m.get(key))
}

// GetStringMapString returns the value of the given key as a map of strings to
// strings or nil.
func (m *KVMem) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(m.get(key))
}

// GetStringMapStringSlice returns the value of the given key as a map of
// strings to string slices or nil.
func (m *KVMem) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(m.get(key))
}

// GetStringSlice returns the value of the given key as a string slice or nil.
func (m *KVMem) GetStringSlice(key string) []string {
	return cast.ToStringSlice(m.get(key))
}

// GetTime returns the value of the given key as a time.
func (m *KVMem) GetTime(key string) time.Time {
	return cast.ToTime(m.get(key))
}

// GetUint returns the value of the given key as an unsigned integer.
func (m *KVMem) GetUint(key string) uint {
	return cast.ToUint(m.get(key))
}

// GetUint16 returns the value of the given key as a 16-bit unsigned integer.
func (m *KVMem) GetUint16(key string) uint16 {
	return cast.ToUint16(m.get(key))
}

// GetUint32 returns the value of the given key as a 32-bit unsigned integer.
func (m *KVMem) GetUint32(key string) uint32 {
	return cast.ToUint32(m.get(key))
}

// GetUint64 returns the value of the given key as a 64-bit unsigned integer.
func (m *KVMem) GetUint64(key string) uint64 {
	return cast.ToUint64(m.get(key))
}

// Sub returns a KVMem pointing to all the same internals as this KVMem, but
// with a prefix set so that all reads ad writes only happen to keys with the
// given prefix.
func (m *KVMem) Sub(key string) KV {
	v := m.get(key)
	if _, isStringMap := v.(map[string]any); isStringMap {
		return &KVMem{
			prefix:       key,
			values:       m.values,
			requirements: m.requirements,
			aliases:      m.aliases,
		}
	}
	return nil
}

// IsSet returns true if the key has been set.
func (m *KVMem) IsSet(key string) bool {
	_, exists := m.getExists(key)
	return exists
}

// Clear resets the storage to empty.
func (m *KVMem) Clear() {
	m.values = make(map[string]any, 10)
}

// Set sets the key to the given value.
func (m *KVMem) Set(key string, value any) {
	m.set(key, value)
}

// Update replaces the top level keys with the given values. This does not merge.
func (m *KVMem) Update(values map[string]any) {
	for k, v := range values {
		m.set(k, v)
	}
}

// UpdateStrings sets all the keys in the map. This effectively does a merge.
func (m *KVMem) UpdateStrings(values map[string]string) {
	for k, v := range values {
		m.set(k, v)
	}
}

// RegisterAlias creates an alias so that any attempt to get or set alias will
// get or set key instead.
func (m *KVMem) RegisterAlias(alias, key string) {
	m.aliases[m.key(alias)] = m.key(key)
}

// MarkRequired adds the given key to the list of required keys.
func (m *KVMem) MarkRequired(key string) {
	m.requirements[m.key(key)] = struct{}{}
}

// IsRequired returns true if the given key has been marked as required.
func (m *KVMem) IsRequired(key string) bool {
	_, exists := m.requirements[m.key(key)]
	return exists
}

// IsPassingRequirements returns true if all the required keys have been set.
func (m *KVMem) IsPassingRequirements() bool {
	for k := range m.requirements {
		if !m.IsSet(k) {
			return false
		}
	}
	return true
}

// MissingRequirements returns all the keys that are missing a value from the
// list of required keys.
func (m *KVMem) MissingRequirements() []string {
	missing := make([]string, 0, len(m.requirements))
	for k := range m.requirements {
		if !m.IsSet(k) {
			missing = append(missing, k)
		}
	}
	sort.Strings(missing)
	return missing
}
