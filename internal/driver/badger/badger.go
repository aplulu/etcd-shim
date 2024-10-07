package badger

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dgraph-io/badger/v4"

	"github.com/aplulu/etcd-shim/internal/driver"
	"github.com/aplulu/etcd-shim/internal/driver/registry"
)

const (
	internalPrefix = "_etcd-shim/"
	revisionKey    = "revision"
)

func init() {
	registry.Register("badger", New)
}

func New(ctx context.Context, log *slog.Logger) (driver.Driver, error) {
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		return nil, fmt.Errorf("badger.New: failed to open badger: %w", err)
	}

	return &badgerDriver{
		log: log,
		db:  db,
	}, nil
}

type badgerDriver struct {
	log *slog.Logger
	db  *badger.DB
}

func (d *badgerDriver) Watch(ctx context.Context, key []byte, startRevision int64) chan *driver.WatchEvent {
	ch := make(chan *driver.WatchEvent)

	go func() {
		defer close(ch)

		err := d.db.View(func(txn *badger.Txn) error {
			opts := badger.DefaultIteratorOptions
			opts.PrefetchSize = 10
			it := txn.NewIterator(opts)
			defer it.Close()

			for it.Seek(key); it.Valid(); it.Next() {
				item := it.Item()
				k := item.Key()
				v, err := item.ValueCopy(nil)
				if err != nil {
					return fmt.Errorf("badgerDriver.Watch: failed to copy value: %w", err)
				}

				ch <- &driver.WatchEvent{
					KV: &driver.KeyValue{
						Key:   k,
						Value: v,
					},
				}
			}

			return nil
		})
		if err != nil {
			d.log.Error("badgerDriver.Watch: failed to view", "error", err)
		}
	}()

	return ch
}

func (d *badgerDriver) Range(ctx context.Context, key []byte, end []byte) ([]driver.KeyValue, int64, error) {
	var results []driver.KeyValue

	err := d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(key); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			v, err := item.ValueCopy(nil)
			if err != nil {
				return fmt.Errorf("badgerDriver.Range: failed to copy value: %w", err)
			}

			results = append(results, driver.KeyValue{
				Key:   k,
				Value: v,
			})

			if end != nil && string(k) == string(end) {
				break
			}
		}

		return nil
	})
	if err != nil {
		return nil, 0, fmt.Errorf("badgerDriver.Range: failed to view: %w", err)
	}

	return results, 0, nil
}

func (d *badgerDriver) Put(ctx context.Context, key []byte, value []byte) (int64, error) {
	if err := d.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	}); err != nil {
		return 0, fmt.Errorf("badgerDriver.Put: failed to update: %w", err)
	}

	revision, err := d.incrementRevision(ctx)
	if err != nil {
		return 0, fmt.Errorf("badgerDriver.Put: failed to increment revision: %w", err)
	}

	return revision, nil
}

func (d *badgerDriver) incrementRevision(ctx context.Context) (int64, error) {
	seq, err := d.db.GetSequence([]byte(internalPrefix+revisionKey), 1)
	if err != nil {
		return 0, fmt.Errorf("badgerDriver.incrementRevision: failed to get sequence: %w", err)
	}
	revision, err := seq.Next()
	if err != nil {
		return 0, fmt.Errorf("badgerDriver.incrementRevision: failed to increment sequence: %w", err)
	}

	return int64(revision), nil
}

func (d *badgerDriver) getRevision(ctx context.Context) (int64, error) {
	var revision int64
	if err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(internalPrefix + revisionKey))
		if err != nil {
			return fmt.Errorf("badgerDriver.getRevision: failed to get revision: %w", err)
		}

		return item.Value(func(val []byte) error {
			revision = int64(val[0])
			return nil
		})
	}); err != nil {
		return 0, fmt.Errorf("badgerDriver.getRevision: failed to view: %w", err)
	}

	return revision, nil
}
