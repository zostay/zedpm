package storage

import (
	"sort"
	"time"
)

// Verify that KVChanges is a KV.
var _ KV = &KVChanges{}

// KVChanges is a key-value store that tracks changes atop the Inner KV. If a
// value is changed, it will not change the Inner value. However, if that value
// is retrieved, the changed value will be returned by this object.
type KVChanges struct {
	changes KV

	// Inner is the KV wrapped by this and will not be modified by KVChanges
	// when Set or other write methods are called.
	Inner KV
}

// WithChangeTracking adds change tracking to inner. The returned KVChanges
// object will not modify the inner object when write methods are called upon
// it, but will store the changes that are written.
func WithChangeTracking(inner KV) *KVChanges {
	return &KVChanges{
		Inner:   inner,
		changes: New(),
	}
}

// AllKeys will return all keys that have been set either on the Inner KV or
// added by making writes to this object.
func (c *KVChanges) AllKeys() []string {
	var (
		innerKeys   = c.Inner.AllKeys()
		changesKeys = c.changes.AllKeys()
		set         = make(map[string]struct{}, len(innerKeys)+len(changesKeys))
	)
	for _, k := range innerKeys {
		set[k] = struct{}{}
	}
	for _, k := range changesKeys {
		set[k] = struct{}{}
	}
	out := keys[struct{}](set)
	sort.Strings(out)
	return out
}

// AllSettings returns the complete map of values. This starts with the values
// in the Inner KV and then layers on local changes so that any key lookup in
// the returned map should have the same value there that it would have in the
// KVChanges object.
func (c *KVChanges) AllSettings() map[string]any {
	var (
		innerKeys   = c.Inner.AllKeys()
		changesKeys = c.changes.AllKeys()
		out         = make(map[string]any, len(innerKeys)+len(changesKeys))
	)
	for _, k := range innerKeys {
		out[k] = c.Inner.Get(k)
	}
	for _, k := range changesKeys {
		out[k] = c.changes.Get(k)
	}
	return out
}

// AllSettingsStrings returns the complete map of values as a flat set a strings
// with any local changes overriding those of the Inner KV.
func (c *KVChanges) AllSettingsStrings() map[string]string {
	keys := c.AllKeys()
	out := make(map[string]string, len(keys))
	for _, k := range keys {
		out[k] = c.GetString(k)
	}
	return out
}

// getc just make writing getters simpler.
func getc[T any](c *KVChanges, key string, getter func(KV, string) T) T {
	if c.changes.IsSet(key) {
		return getter(c.changes, key)
	}
	return getter(c.Inner, key)
}

func (c *KVChanges) Get(key string) any {
	return getc[any](c, key, KV.Get)
}

func (c *KVChanges) GetBool(key string) bool {
	return getc[bool](c, key, KV.GetBool)
}

func (c *KVChanges) GetDuration(key string) time.Duration {
	return getc[time.Duration](c, key, KV.GetDuration)
}

func (c *KVChanges) GetFloat64(key string) float64 {
	return getc[float64](c, key, KV.GetFloat64)
}

func (c *KVChanges) GetInt(key string) int {
	return getc[int](c, key, KV.GetInt)
}

func (c *KVChanges) GetInt32(key string) int32 {
	return getc[int32](c, key, KV.GetInt32)
}

func (c *KVChanges) GetInt64(key string) int64 {
	return getc[int64](c, key, KV.GetInt64)
}

func (c *KVChanges) GetIntSlice(key string) []int {
	return getc[[]int](c, key, KV.GetIntSlice)
}

func (c *KVChanges) GetString(key string) string {
	return getc[string](c, key, KV.GetString)
}

func (c *KVChanges) GetStringMap(key string) map[string]any {
	return getc[map[string]any](c, key, KV.GetStringMap)
}

func (c *KVChanges) GetStringMapString(key string) map[string]string {
	return getc[map[string]string](c, key, KV.GetStringMapString)
}

func (c *KVChanges) GetStringMapStringSlice(key string) map[string][]string {
	return getc[map[string][]string](c, key, KV.GetStringMapStringSlice)
}

func (c *KVChanges) GetStringSlice(key string) []string {
	return getc[[]string](c, key, KV.GetStringSlice)
}

func (c *KVChanges) GetTime(key string) time.Time {
	return getc[time.Time](c, key, KV.GetTime)
}

func (c *KVChanges) GetUint(key string) uint {
	return getc[uint](c, key, KV.GetUint)
}

func (c *KVChanges) GetUint16(key string) uint16 {
	return getc[uint16](c, key, KV.GetUint16)
}

func (c *KVChanges) GetUint32(key string) uint32 {
	return getc[uint32](c, key, KV.GetUint32)
}

func (c *KVChanges) GetUint64(key string) uint64 {
	return getc[uint64](c, key, KV.GetUint64)
}

func (c *KVChanges) Sub(key string) KV {
	return &KVChanges{
		changes: c.changes.Sub(key),
		Inner:   c.Inner.Sub(key),
	}
}

// IsSet returns true if the given key has been set in either the changes or the
// Inner.
func (c *KVChanges) IsSet(key string) bool {
	return c.changes.IsSet(key) || c.Inner.IsSet(key)
}

// TODO Should Clear apply changes to make all the Inner disappear instead without modifying the Inner?

// Clear will clear both the changes and the Inner.
func (c *KVChanges) Clear() {
	c.changes.Clear()
	c.Inner.Clear()
}

// Set adds a setting to local changes and does not modify the Inner.
func (c *KVChanges) Set(key string, value any) {
	c.changes.Set(key, value)
}

// Update applies an update from the given map to local changes and odes not
// modify the Inner.
func (c *KVChanges) Update(values map[string]any) {
	c.changes.Update(values)
}

// UpdateStrings applies an update from the given map to local changes and does
// not modify the Inner.
func (c *KVChanges) UpdateStrings(values map[string]string) {
	c.changes.UpdateStrings(values)
}

// RegisterAlias registers an alias in the local changes and in the Inner.
func (c *KVChanges) RegisterAlias(alias, key string) {
	c.changes.RegisterAlias(alias, key)
	c.Inner.RegisterAlias(alias, key)
}

// Changes returns all the local changes as a hierarchical map of values.
func (c *KVChanges) Changes() map[string]any {
	return c.changes.AllSettings()
}

// ChangesStrings returns all the local changes as a flat map of string values.
func (c *KVChanges) ChangesStrings() map[string]string {
	changesKeys := c.changes.AllKeys()
	out := make(map[string]string, len(changesKeys))
	for _, key := range changesKeys {
		out[key] = c.changes.GetString(key)
	}
	return out
}

// ClearChanges clears the local changes without modifying the Inner.
func (c *KVChanges) ClearChanges() {
	c.changes.Clear()
}
