package main

import (
	"context"
	"github.com/j03hanafi/seternak-backend/logger"
	"github.com/j03hanafi/seternak-backend/router"
	"go.uber.org/zap"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var err error
	l := logger.Get()
	defer func(l *zap.Logger) {
		_ = l.Sync()
	}(l)

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	defer stop()

	app := router.New()

	go func() {
		if err = app.Listen(":8080"); err != nil {
			log.Fatal("Server Error", err)
		}
	}()

	l.Info("Server is starting...")

	<-ctx.Done()
	stop()
	l.Info("shutting down gracefully, press Ctrl+C again to force")

	err = app.ShutdownWithTimeout(5 * time.Second)
	if err != nil {
		l.Fatal("Server forced to shutdown", zap.Error(err))
	}

	l.Info("Server was successful shutdown.")
}
