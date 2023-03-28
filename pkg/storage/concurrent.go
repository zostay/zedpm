package storage

import (
	"sync"
	"time"
)

// Viery that KVCon is a KV.
var _ KV = &KVCon{}

// KVCon wraps a KV in synchronization tooling that prevents concurrent
// modifications. This uses a sync.RWMutex so that many reads can happen
// simultaneously, but writes must be exclusive.
type KVCon struct {
	inner KV
	lock  *sync.RWMutex
}

// WithSafeConcurrency wraps the given KV in a concurrency safe KVCon.
func WithSafeConcurrency(inner KV) *KVCon {
	return WithLock(inner, &sync.RWMutex{})
}

// WithLock wraps the given KV in a concurrency safe KVCon using the given lock.
func WithLock(inner KV, lock *sync.RWMutex) *KVCon {
	return &KVCon{inner, lock}
}

// readSyncUU makes writing getters easier.
func readSyncUU[Out any](
	c *KVCon,
	call func(KV) Out,
) Out {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return call(c.inner)
}

// reaSyncBU makes writing getters easier.
func readSyncBU[In, Out any](
	c *KVCon,
	in In,
	call func(KV, In) Out,
) Out {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return call(c.inner, in)
}

// AllKeys retrieves all keys after acquiring a read lock.
func (c *KVCon) AllKeys() []string {
	return readSyncUU[[]string](c, KV.AllKeys)
}

// AllSettings retrieves all settings after acquiring a read lock.
func (c *KVCon) AllSettings() map[string]any {
	return readSyncUU[map[string]any](c, KV.AllSettings)
}

// AllSettingsStrings gets all settings after acquiring a read lock.
func (c *KVCon) AllSettingsStrings() map[string]string {
	return readSyncUU[map[string]string](c, KV.AllSettingsStrings)
}

func (c *KVCon) Get(key string) any {
	return readSyncBU[string, any](c, key, KV.Get)
}

func (c *KVCon) GetBool(key string) bool {
	return readSyncBU[string, bool](c, key, KV.GetBool)
}

func (c *KVCon) GetDuration(key string) time.Duration {
	return readSyncBU[string, time.Duration](c, key, KV.GetDuration)
}

func (c *KVCon) GetFloat64(key string) float64 {
	return readSyncBU[string, float64](c, key, KV.GetFloat64)
}

func (c *KVCon) GetInt(key string) int {
	return readSyncBU[string, int](c, key, KV.GetInt)
}

func (c *KVCon) GetInt32(key string) int32 {
	return readSyncBU[string, int32](c, key, KV.GetInt32)
}

func (c *KVCon) GetInt64(key string) int64 {
	return readSyncBU[string, int64](c, key, KV.GetInt64)
}

func (c *KVCon) GetIntSlice(key string) []int {
	return readSyncBU[string, []int](c, key, KV.GetIntSlice)
}

func (c *KVCon) GetString(key string) string {
	return readSyncBU[string, string](c, key, KV.GetString)
}

func (c *KVCon) GetStringMap(key string) map[string]any {
	return readSyncBU[string, map[string]any](c, key, KV.GetStringMap)
}

func (c *KVCon) GetStringMapString(key string) map[string]string {
	return readSyncBU[string, map[string]string](c, key, KV.GetStringMapString)
}

func (c *KVCon) GetStringMapStringSlice(key string) map[string][]string {
	return readSyncBU[string, map[string][]string](c, key, KV.GetStringMapStringSlice)
}

func (c *KVCon) GetStringSlice(key string) []string {
	return readSyncBU[string, []string](c, key, KV.GetStringSlice)
}

func (c *KVCon) GetTime(key string) time.Time {
	return readSyncBU[string, time.Time](c, key, KV.GetTime)
}

func (c *KVCon) GetUint(key string) uint {
	return readSyncBU[string, uint](c, key, KV.GetUint)
}

func (c *KVCon) GetUint16(key string) uint16 {
	return readSyncBU[string, uint16](c, key, KV.GetUint16)
}

func (c *KVCon) GetUint32(key string) uint32 {
	return readSyncBU[string, uint32](c, key, KV.GetUint32)
}

func (c *KVCon) GetUint64(key string) uint64 {
	return readSyncBU[string, uint64](c, key, KV.GetUint64)
}

// Sub returns an object that retrieves a subset of the values guarded by the
// same lock.
func (c *KVCon) Sub(key string) KV {
	return &KVCon{c.inner.Sub(key), c.lock}
}

// IsSet reports whether this value is set after acquiring a read lock.
func (c *KVCon) IsSet(key string) bool {
	return readSyncBU[string, bool](c, key, KV.IsSet)
}

// Clear deletes all value in the inner KV after acquiring a write lock.
func (c *KVCon) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.inner.Clear()
}

// Set sets a value on the inner KV after acquiring a write lock.
func (c *KVCon) Set(key string, value any) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.inner.Set(key, value)
}

// Update applies an update after acquiring a write lock.
func (c *KVCon) Update(values map[string]any) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.inner.Update(values)
}

// UpdateStrings applies an update after acquiring a write lock.
func (c *KVCon) UpdateStrings(values map[string]string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.inner.UpdateStrings(values)
}

// RegisterAlias adds an alias to the inner KV after acquiring a write lock.
func (c *KVCon) RegisterAlias(alias, key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.inner.RegisterAlias(alias, key)
}

// Atomic performs the given functional call after acquiring a write lock. The
// callback will receive the inner KV as an argument. This should be used for
// any operation that needs to be performed on the inner KV that would have a
// race condition if performed without a mutex lock around the whole operation.
func (c *KVCon) Atomic(call func(KV)) {
	c.lock.Lock()
	defer c.lock.Unlock()
	call(c.inner)
}
