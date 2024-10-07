package main

import (
	"context"
	"fmt"
	"github.com/aplulu/etcd-shim/internal/server"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aplulu/etcd-shim/internal/config"
)

func main() {
	if err := config.LoadConf(); err != nil {
		panic(err)
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-quitCh
		log.Info("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.StopServer(shutdownCtx); err != nil {
			log.Error(fmt.Sprintf("command.ServeCommand: failed to stop server: %+v", err))
			os.Exit(1)
			return
		}
	}()

	log.Info("Starting server...")
	if err := server.StartServer(log); err != nil {
		log.Error(fmt.Sprintf("command.ServeCommand: failed to start server: %+v", err))
		os.Exit(1)
	}
}
