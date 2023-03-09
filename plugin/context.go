package plugin

import (
	"context"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"

	"github.com/zostay/zedpm/config"
	"github.com/zostay/zedpm/pkg/storage"
)

type contextKey struct{}

// Context is a container providing the information accessed by the various
// context.Context accessors used by plugins.
//
// This object is safe to use concurrently. Changes to properties are
// synchronized as are cleanup functions and add files. However, if your changes
// are not idempotent, special care must be taken to avoid race conditions. The
// AtomicProperties function is provided to allow for atomic operations to be
// safely performed without races.
type Context struct {
	// TODO The hclog.Logger is rubbish. Switch to zap.Logger.
	logger     hclog.Logger
	cleanup    []SimpleTask
	addFiles   []string
	properties *storage.KVChanges
	safeProps  *storage.KVCon
	lock       *sync.Mutex
}

// SimpleTask is the type of function used for cleanup functions.
type SimpleTask func()

// NewContext constructs a new plugin Context, ready to be used with
// InitializeContext to attach it to a context.Context.
//
// The properties object here will not be changed when any of the context
// setters, such as Set or ApplyChanges are used. These will be stored by
// wrapping the given properties in a storage.KVChanges object. The original
// caller can apply these changes by using UpdateStorage and StorageChanges
// together.
func NewContext(
	logger hclog.Logger,
	properties storage.KV,
) *Context {
	props := storage.WithChangeTracking(properties)
	return &Context{
		logger:     logger,
		cleanup:    make([]SimpleTask, 0, 10),
		addFiles:   make([]string, 0, 10),
		properties: props,
		safeProps:  storage.WithSafeConcurrency(props),
		lock:       &sync.Mutex{},
	}
}

// NewConfigContext constructs a new plugin Context, but with the given scopes
// defined for constructing properties to provide. This is ready to be used
// with InitializeContext to attach it to a context.Context.
//
// This has the same limitations as NewContext regarding the properties that is
// mentioned in NewContext.
func NewConfigContext(
	logger hclog.Logger,
	runtime storage.KV,
	taskName string,
	targetName string,
	pluginName string,
	cfg *config.Config,
) *Context {
	return NewContext(logger, cfg.ToKV(runtime, taskName, targetName, pluginName))
}

// UpdateStorage allows the owner of the Context to update the properties of the
// properties object given during construction. None of the mutator functions
// that work on context.Context are able to make changes to the original
// properties object as they only apply their changes to a storage.KVChanges
// wrapper applied during Context construction.
func (p *Context) UpdateStorage(store map[string]string) {
	// This is a bit odd, but hear me out... the safeProps lock is providing us
	// with a write-safe mutex here. We could write the following as:
	//
	//   p.safeProps.Atomic(func(kv storage.KV) {
	//       kv.(*KVChanges).Inner.UpdateStrings(store)
	//   })
	//
	// But we already have "kv.(*storage.KVChanges)" in p.properties, so let's
	// not worry about the type coercion thing. However, if ever p.properties !=
	// kv(*storage.KVChanges), this will (probably) break in an ugly way.
	p.safeProps.Atomic(func(storage.KV) {
		p.properties.Inner.UpdateStrings(store)
	})
}

// StorageChanges clears any changes that were made by callers to the mutator
// methods on the context.Context and returns them. These can be made permanent
// by calling UpdateStorage.
func (p *Context) StorageChanges() map[string]string {
	var changes map[string]string
	// See the comment in UpdateStorage regarding why this is written this
	// way...
	p.safeProps.Atomic(func(storage.KV) {
		changes = p.properties.ChangesStrings()
		p.properties.ClearChanges()
	})
	return changes
}

// InitializeContext attaches the plugin.Context to the context.Context.
func InitializeContext(ctx context.Context, pctx *Context) context.Context {
	return context.WithValue(ctx, contextKey{}, pctx)
}

// contextFrom is the internal method used to extract the plugin Context from
// the context.Context.
func contextFrom(ctx context.Context) *Context {
	v := ctx.Value(contextKey{})
	pctx, isPluginContext := v.(*Context)
	if !isPluginContext {
		panic("context is missing plugin configuration")
	}
	return pctx
}

// Logger returns the logger for the plugin to use. If any withArgs are passed,
// the returned logger will have had the hclog.Logger.With function called to
// set properties on the logger.
func Logger(ctx context.Context, withArgs ...interface{}) hclog.Logger {
	pctx := contextFrom(ctx)
	if len(withArgs) > 0 {
		return pctx.logger.With(withArgs...)
	}
	return pctx.logger
}

// ForCleanup adds the given task to be performed at cleanup time.
func ForCleanup(ctx context.Context, newCleaner SimpleTask) {
	pctx := contextFrom(ctx)
	pctx.lock.Lock()
	defer pctx.lock.Unlock()
	pctx.cleanup = append(pctx.cleanup, newCleaner)
}

// ListCleanupTasks returns all the cleanup tasks that have been setup so far
// since the start of this Context. The tasks are returned in the reverse order
// they were added.
func ListCleanupTasks(ctx context.Context) []SimpleTask {
	// TODO Implement usage of cleanup tasks.
	// TODO ListCleanupTasks should probably be moved up onto pctx itself.
	pctx := contextFrom(ctx)
	pctx.lock.Lock()
	defer pctx.lock.Unlock()
	tasks := make([]SimpleTask, len(pctx.cleanup))
	for i, f := range pctx.cleanup {
		tasks[len(tasks)-i-1] = f
	}
	return tasks
}

// ToAdd names a new file that has been created or added by the tooling, which
// allows it to be added to VC and other such tools.
func ToAdd(ctx context.Context, newFile string) {
	pctx := contextFrom(ctx)
	pctx.lock.Lock()
	defer pctx.lock.Unlock()
	pctx.addFiles = append(pctx.addFiles, newFile)
}

// ListAdded returns all the files that have been created or added by the
// tooling since the Context was created.
func ListAdded(ctx context.Context) []string {
	pctx := contextFrom(ctx)
	pctx.lock.Lock()
	defer pctx.lock.Unlock()
	return pctx.addFiles
}

// AtomicProperties executes the given function inside a write lock on the
// storage object. This allows non-idempotent, atomic modifications to be made
// to the storage without races.
//
// For example, NEVER RUN:
//
//	// NEVER DO THIS! This contains a race condition and will work inconsistently.
//	v := plugin.GetInt(ctx, "foo")
//	plugin.Set(ctx, "foo", v+1)
//
// Instead, you should follow this example:
//
//	plugin.AtomicProperties(ctx, func(p storage.KV) {
//	    v := p.GetInt("foo")
//	    p.Set("foo", v+1)
//	})
func AtomicProperties(ctx context.Context, atomicOp func(storage.KV)) {
	pctx := contextFrom(ctx)
	pctx.safeProps.Atomic(atomicOp)
}

// IsSet returns whether or not the given key is set in the Context properties.
func IsSet(ctx context.Context, key string) bool {
	pctx := contextFrom(ctx)
	return pctx.safeProps.IsSet(key)
}

// Set sets the given key/value pair on the Context properties.
func Set(ctx context.Context, key string, value any) {
	pctx := contextFrom(ctx)
	pctx.safeProps.Set(key, value)
}

// get[T any] is an internal function that makes writing getters much simpler.
func get[T any](ctx context.Context, key string, getter func(storage.KV, string) T) T {
	pctx := contextFrom(ctx)
	return getter(pctx.safeProps, key)
}

// Get returns the underlying value as is without attempt to coerce the value in
// any way.
func Get(ctx context.Context, key string) any {
	return get(ctx, key, storage.KV.Get)
}

// GetBool returns the value for the given key as a boolean.
func GetBool(ctx context.Context, key string) bool {
	return get(ctx, key, storage.KV.GetBool)
}

// GetDuration returns the value for the given key as a time.Duration.
func GetDuration(ctx context.Context, key string) time.Duration {
	return get(ctx, key, storage.KV.GetDuration)
}

// GetFloat64 returns the value for the given key as a float64.
func GetFloat64(ctx context.Context, key string) float64 {
	return get(ctx, key, storage.KV.GetFloat64)
}

// GetInt returns the value for the given key as an int.
func GetInt(ctx context.Context, key string) int {
	return get(ctx, key, storage.KV.GetInt)
}

// GetInt32 returns the value for the given key as an int32.
func GetInt32(ctx context.Context, key string) int32 {
	return get(ctx, key, storage.KV.GetInt32)
}

// GetInt64 returns the value for the given key as an int64.
func GetInt64(ctx context.Context, key string) int64 {
	return get(ctx, key, storage.KV.GetInt64)
}

// GetIntSlice returns the value for the given key as an []int.
func GetIntSlice(ctx context.Context, key string) []int {
	return get(ctx, key, storage.KV.GetIntSlice)
}

// GetString returns the value for the given key as a string.
func GetString(ctx context.Context, key string) string {
	return get(ctx, key, storage.KV.GetString)
}

// GetStringMap returns the value for the given key as a map[string]any.
func GetStringMap(ctx context.Context, key string) map[string]any {
	return get(ctx, key, storage.KV.GetStringMap)
}

// GetStringMapString returns the value for the given key as a map[string]string.
func GetStringMapString(ctx context.Context, key string) map[string]string {
	return get(ctx, key, storage.KV.GetStringMapString)
}

// GetStringMapStringSlice returns the value for the given key as a
// map[string][]string.
func GetStringMapStringSlice(ctx context.Context, key string) map[string][]string {
	return get(ctx, key, storage.KV.GetStringMapStringSlice)
}

// GetStringSlice returns the value for the given key as a []string.
func GetStringSlice(ctx context.Context, key string) []string {
	return get(ctx, key, storage.KV.GetStringSlice)
}

// GetTime returns the value for the given key as a time.Time.
func GetTime(ctx context.Context, key string) time.Time {
	return get(ctx, key, storage.KV.GetTime)
}

// GetUint returns the value for the given key as a uint.
func GetUint(ctx context.Context, key string) uint {
	return get(ctx, key, storage.KV.GetUint)
}

// GetUint16 returns the value for the given key as a uint16.
func GetUint16(ctx context.Context, key string) uint16 {
	return get(ctx, key, storage.KV.GetUint16)
}

// GetUint32 returns the value for the given key as a uint32.
func GetUint32(ctx context.Context, key string) uint32 {
	return get(ctx, key, storage.KV.GetUint32)
}

// GetUint64 returns the value for the given key as a uint64.
func GetUint64(ctx context.Context, key string) uint64 {
	return get(ctx, key, storage.KV.GetUint64)
}

// KV returns the properties attached to the context.
func KV(ctx context.Context) *storage.KVCon {
	pctx := contextFrom(ctx)
	return pctx.safeProps
}

// ApplyChanges will apply the given changes to the properties.
func ApplyChanges(ctx context.Context, changes map[string]string) {
	pctx := contextFrom(ctx)
	pctx.safeProps.UpdateStrings(changes)
}
