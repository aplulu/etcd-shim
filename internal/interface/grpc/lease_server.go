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

type leaseServer struct {
	log *slog.Logger
}

func (s *leaseServer) LeaseGrant(ctx context.Context, req *etcdserverpb.LeaseGrantRequest) (*etcdserverpb.LeaseGrantResponse, error) {
	return &etcdserverpb.LeaseGrantResponse{
		Header: &etcdserverpb.ResponseHeader{},
		ID:     req.ID,
		TTL:    req.TTL,
	}, nil
}

func (s *leaseServer) LeaseRevoke(ctx context.Context, req *etcdserverpb.LeaseRevokeRequest) (*etcdserverpb.LeaseRevokeResponse, error) {
	return nil, fmt.Errorf("grpc.LeaseRevoke is not implemented: %w", ErrNotImplemented)
}

func (s *leaseServer) LeaseKeepAlive(server etcdserverpb.Lease_LeaseKeepAliveServer) error {
	return fmt.Errorf("grpc.LeaseKeepAlive is not implemented: %w", ErrNotImplemented)
}

func (s *leaseServer) LeaseTimeToLive(ctx context.Context, req *etcdserverpb.LeaseTimeToLiveRequest) (*etcdserverpb.LeaseTimeToLiveResponse, error) {
	return nil, fmt.Errorf("grpc.LeaseTimeToLive is not implemented: %w", ErrNotImplemented)
}

func (s *leaseServer) LeaseLeases(ctx context.Context, req *etcdserverpb.LeaseLeasesRequest) (*etcdserverpb.LeaseLeasesResponse, error) {
	return nil, fmt.Errorf("grpc.LeaseLeases is not implemented: %w", ErrNotImplemented)
}

func RegisterLeaseServer(ctx context.Context, gs *grpc.Server, mux *runtime.ServeMux, l *slog.Logger) error {
	s := &leaseServer{
		log: l,
	}
	etcdserverpb.RegisterLeaseServer(gs, s)
	if err := gw.RegisterLeaseHandlerServer(ctx, mux, s); err != nil {
		return fmt.Errorf("failed to register ClusterServer: %w", err)
	}

	return nil
}
