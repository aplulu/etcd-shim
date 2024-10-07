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

type clusterServer struct {
	log *slog.Logger
}

func (s *clusterServer) MemberAdd(ctx context.Context, req *etcdserverpb.MemberAddRequest) (*etcdserverpb.MemberAddResponse, error) {
	return nil, fmt.Errorf("not implemented: MemberAdd: %w", ErrNotImplemented)
}

func (s *clusterServer) MemberRemove(ctx context.Context, req *etcdserverpb.MemberRemoveRequest) (*etcdserverpb.MemberRemoveResponse, error) {
	return nil, fmt.Errorf("not implemented: MemberRemove: %w", ErrNotImplemented)
}

func (s *clusterServer) MemberUpdate(ctx context.Context, req *etcdserverpb.MemberUpdateRequest) (*etcdserverpb.MemberUpdateResponse, error) {
	return nil, fmt.Errorf("not implemented: MemberUpdate: %w", ErrNotImplemented)
}

func (s *clusterServer) MemberList(ctx context.Context, req *etcdserverpb.MemberListRequest) (*etcdserverpb.MemberListResponse, error) {
	return &etcdserverpb.MemberListResponse{
		Header: &etcdserverpb.ResponseHeader{},
		Members: []*etcdserverpb.Member{
			{
				ID:         0,
				Name:       "etcd-shim",
				PeerURLs:   nil,
				ClientURLs: nil,
				IsLearner:  false,
			},
		},
	}, nil
}

func (s *clusterServer) MemberPromote(ctx context.Context, req *etcdserverpb.MemberPromoteRequest) (*etcdserverpb.MemberPromoteResponse, error) {
	return nil, fmt.Errorf("not implemented: MemberPromote: %w", ErrNotImplemented)
}

func RegisterClusterServer(ctx context.Context, gs *grpc.Server, mux *runtime.ServeMux, l *slog.Logger) error {
	s := &clusterServer{
		log: l,
	}
	etcdserverpb.RegisterClusterServer(gs, s)
	if err := gw.RegisterClusterHandlerServer(ctx, mux, s); err != nil {
		return fmt.Errorf("failed to register ClusterServer: %w", err)
	}

	return nil
}
