package grpc

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.etcd.io/etcd/api/v3/etcdserverpb"
	"go.etcd.io/etcd/api/v3/etcdserverpb/gw"
	"google.golang.org/grpc"
)

type maintenanceServer struct {
	log *slog.Logger
}

func (s *maintenanceServer) Alarm(ctx context.Context, request *etcdserverpb.AlarmRequest) (*etcdserverpb.AlarmResponse, error) {
	return nil, fmt.Errorf("not implemented: Alarm: %w", ErrNotImplemented)
}

func (s *maintenanceServer) Status(ctx context.Context, request *etcdserverpb.StatusRequest) (*etcdserverpb.StatusResponse, error) {
	s.log.Info(
		"Status",
	)

	return &etcdserverpb.StatusResponse{
		Header:  &etcdserverpb.ResponseHeader{},
		DbSize:  0,
		Version: "3.5.0",
	}, nil
}

func (s *maintenanceServer) Defragment(ctx context.Context, request *etcdserverpb.DefragmentRequest) (*etcdserverpb.DefragmentResponse, error) {
	return nil, fmt.Errorf("not implemented: Defragment: %w", ErrNotImplemented)
}

func (s *maintenanceServer) Hash(ctx context.Context, request *etcdserverpb.HashRequest) (*etcdserverpb.HashResponse, error) {
	return nil, fmt.Errorf("not implemented: Hash: %w", ErrNotImplemented)
}

func (s *maintenanceServer) HashKV(ctx context.Context, request *etcdserverpb.HashKVRequest) (*etcdserverpb.HashKVResponse, error) {
	return nil, fmt.Errorf("not implemented: HashKV: %w", ErrNotImplemented)
}

func (s *maintenanceServer) Snapshot(request *etcdserverpb.SnapshotRequest, server etcdserverpb.Maintenance_SnapshotServer) error {
	return fmt.Errorf("not implemented: Snapshot: %w", ErrNotImplemented)
}

func (s *maintenanceServer) MoveLeader(ctx context.Context, request *etcdserverpb.MoveLeaderRequest) (*etcdserverpb.MoveLeaderResponse, error) {
	return nil, fmt.Errorf("not implemented: MoveLeader: %w", ErrNotImplemented)
}

func (s *maintenanceServer) Downgrade(ctx context.Context, request *etcdserverpb.DowngradeRequest) (*etcdserverpb.DowngradeResponse, error) {
	return nil, fmt.Errorf("not implemented: Downgrade: %w", ErrNotImplemented)
}

func RegisterMaintenanceServer(ctx context.Context, gs *grpc.Server, mux *runtime.ServeMux, l *slog.Logger) error {
	s := &maintenanceServer{
		log: l,
	}
	etcdserverpb.RegisterMaintenanceServer(gs, s)
	if err := gw.RegisterMaintenanceHandlerServer(ctx, mux, s); err != nil {
		return fmt.Errorf("failed to register maintenance server: %w", err)
	}

	return nil
}
