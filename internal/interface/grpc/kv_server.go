package grpc

import (
	"context"
	"fmt"
	"github.com/aplulu/etcd-shim/internal/driver"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"log/slog"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.etcd.io/etcd/api/v3/etcdserverpb"
	"go.etcd.io/etcd/api/v3/etcdserverpb/gw"
	"google.golang.org/grpc"
)

type kvServer struct {
	log    *slog.Logger
	driver driver.Driver
}

func (s *kvServer) Range(ctx context.Context, req *etcdserverpb.RangeRequest) (*etcdserverpb.RangeResponse, error) {
	s.log.Info(
		"Range",
		"key", string(req.Key),
		"range_end", string(req.RangeEnd),
	)

	results, err := s.driver.Range(ctx, req.Key, req.RangeEnd)
	if err != nil {
		s.log.Error("failed to range", "error", err)
		return nil, fmt.Errorf("failed to range: %w", err)
	}

	kvs := make([]*mvccpb.KeyValue, len(results))
	for i, r := range results {
		kvs[i] = &mvccpb.KeyValue{
			Key:     r.Key,
			Value:   r.Value,
			Version: r.Version,
		}
	}

	return &etcdserverpb.RangeResponse{
		Header: &etcdserverpb.ResponseHeader{
			Revision: 1,
		},
		Kvs: nil,
	}, nil
}

func (s *kvServer) Put(ctx context.Context, req *etcdserverpb.PutRequest) (*etcdserverpb.PutResponse, error) {
	s.log.Info(
		"Put",
		"key", string(req.Key),
		"value", string(req.Value),
		"ignore_value", req.IgnoreValue,
		"ignore_lease", req.IgnoreLease,
		"lease", req.Lease,
		"prev_kv", req.PrevKv,
	)

	revision, err := s.driver.Put(ctx, req.Key, req.Value)
	if err != nil {
		s.log.Error("failed to put", "error", err)
		return nil, fmt.Errorf("failed to put: %w", err)
	}

	return &etcdserverpb.PutResponse{
		Header: &etcdserverpb.ResponseHeader{
			Revision: revision,
		},
		PrevKv: nil,
	}, nil
}

func (s *kvServer) DeleteRange(ctx context.Context, req *etcdserverpb.DeleteRangeRequest) (*etcdserverpb.DeleteRangeResponse, error) {
	s.log.Info(
		"DeleteRange",
		"key", string(req.Key),
		"range_end", string(req.RangeEnd),
	)
	return &etcdserverpb.DeleteRangeResponse{}, nil
}

func (s *kvServer) Txn(ctx context.Context, req *etcdserverpb.TxnRequest) (*etcdserverpb.TxnResponse, error) {
	s.log.Info(
		"Txn",
		"compare", req.Compare,
		"success", req.Success,
	)
	return &etcdserverpb.TxnResponse{}, nil
}

func (s *kvServer) Compact(ctx context.Context, req *etcdserverpb.CompactionRequest) (*etcdserverpb.CompactionResponse, error) {
	s.log.Info(
		"Compact",
		"revision", req.Revision,
		"physical", req.Physical,
	)

	return &etcdserverpb.CompactionResponse{}, nil
}

func RegisterKV(ctx context.Context, gs *grpc.Server, mux *runtime.ServeMux, l *slog.Logger, drv driver.Driver) error {
	s := &kvServer{
		log:    l,
		driver: drv,
	}
	etcdserverpb.RegisterKVServer(gs, s)
	if err := gw.RegisterKVHandlerServer(ctx, mux, s); err != nil {
		return fmt.Errorf("failed to register KVServer: %w", err)
	}

	return nil
}
