package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"

	"github.com/aplulu/etcd-shim/internal/config"
	_ "github.com/aplulu/etcd-shim/internal/driver/badger"
	"github.com/aplulu/etcd-shim/internal/driver/registry"
	interfacegrpc "github.com/aplulu/etcd-shim/internal/interface/grpc"
)

var server http.Server

func StartServer(log *slog.Logger) error {
	ctx := context.Background()

	drv, err := registry.NewDriver(config.Driver(), ctx, log)
	if err != nil {
		return fmt.Errorf("failed to create driver: %w", err)
	}

	grpcServer := grpc.NewServer()
	gwMux := runtime.NewServeMux()

	if err := interfacegrpc.RegisterKV(ctx, grpcServer, gwMux, log, drv); err != nil {
		return fmt.Errorf("server.StartServer: failed to register KVServer: %w", err)
	}
	if err := interfacegrpc.RegisterWatch(ctx, grpcServer, gwMux, log); err != nil {
		return fmt.Errorf("server.StartServer: failed to register WatchServer: %w", err)
	}
	if err := interfacegrpc.RegisterClusterServer(ctx, grpcServer, gwMux, log); err != nil {
		return fmt.Errorf("server.StartServer: failed to register ClusterServer: %w", err)
	}
	if err := interfacegrpc.RegisterMaintenanceServer(ctx, grpcServer, gwMux, log); err != nil {
		return fmt.Errorf("server.StartServer: failed to register maintenance server: %w", err)
	}
	if err := interfacegrpc.RegisterLeaseServer(ctx, grpcServer, gwMux, log); err != nil {
		return fmt.Errorf("server.StartServer: failed to register LeaseServer: %w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("ok")); err != nil {
			log.ErrorContext(r.Context(), fmt.Sprintf("server.StartServer: failed to write response: %s", err.Error()))
		}
	})
	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		res := struct {
			ETCDServer  string `json:"etcdserver"`
			ETCDCluster string `json:"etcdcluster"`
		}{
			ETCDServer:  config.ETCDVersion(),
			ETCDCluster: config.ETCDClusterVersion(),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.ErrorContext(r.Context(), fmt.Sprintf("server.StartServer: failed to write response: %s", err.Error()))
		}
	})

	server = http.Server{
		Addr: net.JoinHostPort(config.Listen(), config.Port()),
		Handler: h2c.NewHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Info(
					"Request",
					"method", r.Method,
					"path", r.URL.Path,
					"proto", r.Proto,
					"remote", r.RemoteAddr,
					"content-type", r.Header.Get("Content-Type"),
				)

				if r.ProtoMajor == 2 && r.Header.Get("Content-Type") == "application/grpc" {
					grpcServer.ServeHTTP(w, r)
				} else if strings.HasPrefix(r.URL.Path, "/v3") {
					gwMux.ServeHTTP(w, r)
				} else {
					mux.ServeHTTP(w, r)
				}
			}),
			&http2.Server{},
		),
	}

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func StopServer(ctx context.Context) error {
	return server.Shutdown(ctx)
}
