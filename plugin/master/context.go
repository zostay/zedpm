package master

import (
	"context"
	"sort"
	"sync"

	"github.com/zostay/zedpm/pkg/storage"
	"github.com/zostay/zedpm/plugin/grpc/client"
)

var _ client.Context = &PluginTaskContext{}

// PhaseContext provides a base context for use in tracking the state related to
// an execution phase.
type PhaseContext struct {
	properties *storage.KVChanges  // changes to properties during this phase
	phaseFiles map[string]struct{} // files added by tasks in this phase so far
	lock       sync.RWMutex        // this lock keeps this context synchronized
}

// PluginTaskContext provides a value to be stored in context.Context to track
// the state of execution for the master interface and executor for a particular
// phase, task, target, and plugin. This object is accessible to the plugins.
type PluginTaskContext struct {
	*PhaseContext
	configProps  storage.KV // properties built from the configuration for the current task, target, phase, etc.
	localChanges storage.KV // changes to properties from previous phases
}

// NewContext constructs and returns a new phase context.
func NewContext(properties storage.KV) *PhaseContext {
	return &PhaseContext{
		properties: storage.WithChangeTracking(properties),
		phaseFiles: make(map[string]struct{}, 10),
	}
}

// withPluginTask associates a new plugin/task context with the given context
// for use with the gRPC client context interface.
func (pc *PhaseContext) withPluginTask(
	ctx context.Context,
	configProps storage.KV,
) context.Context {
	return client.WithContext(ctx, &PluginTaskContext{
		PhaseContext: pc,
		configProps:  configProps,
		localChanges: pc.properties.Inner,
	})
}

// KV returns the configuration properties for the current task and plugin along
// with any per-phase changes that have been accumulated thus far.
func (ptc *PluginTaskContext) KV() *storage.KVCon {
	return storage.WithLock(
		storage.Layers(
			ptc.localChanges,
			ptc.configProps,
		), &ptc.lock,
	)
}

// nextPhase transitions a phase context to the next phase by absorbing all the
// changes from associated plugin/task contexts. It then resets teh plugin task
// list to empty.
func (pc *PhaseContext) nextPhase() {
	pc.properties.Inner.Update(pc.properties.Changes())
	pc.phaseFiles = make(map[string]struct{}, 10)
}

// ApplyChanges safely updates the changes applied to the current phase.
func (pc *PhaseContext) ApplyChanges(changes map[string]string) {
	pc.lock.Lock()
	defer pc.lock.Unlock()
	pc.properties.UpdateStrings(changes)
}

// ListAdded returns the list of files added so far to this phase.
func (pc *PhaseContext) ListAdded() []string {
	pc.lock.RLock()
	defer pc.lock.RUnlock()
	out := make([]string, 0, len(pc.phaseFiles))
	for key := range pc.phaseFiles {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}

// ToAdd adds more files to the phase.
func (pc *PhaseContext) ToAdd(files []string) {
	pc.lock.Lock()
	defer pc.lock.Unlock()
	for _, file := range files {
		pc.phaseFiles[file] = struct{}{}
	}
}
