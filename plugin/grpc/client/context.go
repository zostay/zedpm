package client

import (
	"context"

	"github.com/zostay/zedpm/pkg/storage"
)

type ContextKey struct{}

type Context interface {
	KV() *storage.KVCon
	ApplyChanges(map[string]string)
	ListAdded() []string
	ToAdd([]string)
}

func clientContext(ctx context.Context) Context {
	return ctx.Value(ContextKey{}).(Context)
}

func KV(ctx context.Context) *storage.KVCon {
	return clientContext(ctx).KV()
}

func ApplyChanges(ctx context.Context, changes map[string]string) {
	clientContext(ctx).ApplyChanges(changes)
}

func ListAdded(ctx context.Context) []string {
	return clientContext(ctx).ListAdded()
}

func ToAdd(ctx context.Context, files []string) {
	clientContext(ctx).ToAdd(files)
}
