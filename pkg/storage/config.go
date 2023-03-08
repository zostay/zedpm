package storage

import "errors"

var errCfg = errors.New("configuration is not writable")

// Verify that KVCfg is a KV.
var _ KV = &KVCfg{}

// KVCfg is a wrapper around a KVMem that panics if writes are attempted.
type KVCfg struct {
	KVMem
}

// Clear panics.
func (c *KVCfg) Clear() {
	panic(errCfg)
}

// Set panics.
func (c *KVCfg) Set(string, any) {
	panic(errCfg)
}

// Update panics.
func (c *KVCfg) Update(map[string]any) {
	panic(errCfg)
}

// UpdateStrings panics.
func (c *KVCfg) UpdateStrings(map[string]string) {
	panic(errCfg)
}
