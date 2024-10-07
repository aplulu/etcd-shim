package driver

import "context"

type Driver interface {
	Put(ctx context.Context, key []byte, value []byte) (int64, error)
	Range(ctx context.Context, key []byte, end []byte) ([]KeyValue, int64, error)
	Watch(ctx context.Context, key []byte, revision int64) chan *WatchEvent
}

type KeyValue struct {
	Key     []byte
	Value   []byte
	Version int64
}

type WatchEvent struct {
	KV      *KeyValue
	PrevKV  *KeyValue
	Deleted bool
	Created bool
}
