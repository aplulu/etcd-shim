package grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.etcd.io/etcd/api/v3/etcdserverpb"
	"go.etcd.io/etcd/api/v3/etcdserverpb/gw"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	"google.golang.org/grpc"
)

type watchServer struct {
	log *slog.Logger
}

func (s *watchServer) Watch(server etcdserverpb.Watch_WatchServer) (err error) {
	s.log.Info("Watch")
	w := newWatcher(server)

	errCh := make(chan error, 1)

	go func() {
		if err := w.Start(); err != nil {
			errCh <- err
		}
	}()

	select {
	case err = <-errCh:
		if errors.Is(err, context.Canceled) {
			err = rpctypes.ErrGRPCWatchCanceled
		}
	case <-server.Context().Done():
		err = server.Context().Err()
		if errors.Is(err, context.Canceled) {
			err = rpctypes.ErrGRPCWatchCanceled
		}
	}
	w.Stop()

	return err
}

func RegisterWatch(ctx context.Context, gs *grpc.Server, mux *runtime.ServeMux, l *slog.Logger) error {
	s := &watchServer{
		log: l,
	}
	etcdserverpb.RegisterWatchServer(gs, s)
	if err := gw.RegisterWatchHandlerServer(ctx, mux, s); err != nil {
		return fmt.Errorf("failed to register WatchServer: %w", err)
	}

	return nil
}

func newWatcher(server etcdserverpb.Watch_WatchServer) *watcher {
	return &watcher{
		mutex:       sync.Mutex{},
		watchServer: server,
	}
}

type watcher struct {
	mutex       sync.Mutex
	watchServer etcdserverpb.Watch_WatchServer
}

func (w *watcher) Start() error {
	for {
		req, err := w.watchServer.Recv()
		if err != nil {
			return err
		}

		if req.GetCreateRequest() != nil {
			w.handleCreateRequest(req.GetCreateRequest())
		}
	}
}

func (w *watcher) Stop() {
}

func (w *watcher) handleCreateRequest(req *etcdserverpb.WatchCreateRequest) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

}
