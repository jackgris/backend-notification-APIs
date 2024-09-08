package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackgris/backend-notification-APIs/internal/lognotifier"
	"github.com/jackgris/backend-notification-APIs/internal/notification"
	"github.com/jackgris/backend-notification-APIs/internal/store/userdb"
	"github.com/jackgris/backend-notification-APIs/pkg/logs"
	"github.com/jackgris/backend-notification-APIs/pkg/middleware"
)

func main() {
	ctx := context.Background()
	log := logs.New(os.Stdout)

	err := run(ctx, log)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("Error server shutdown: %s\n", err))
	}
}

func run(ctx context.Context, log *logs.Logger) error {

	db := connectDB(ctx, log)
	defer db.Close(ctx)

	lognotifier := lognotifier.NewLogs(db)
	store := userdb.NewStore(db)
	n := notification.NewNotification(store, lognotifier, log)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /notify", middleware.LogResponse(n.NotifyUsers, log))

	portEnv := os.Getenv("PORT")

	port, err := strconv.Atoi(portEnv)
	if err != nil {
		log.Error(ctx, "Environment variable PORT converting to integer")
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:    ":" + portEnv,
		Handler: mux,
	}

	serverErrors := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0), "Server started at port", port)

		serverErrors <- srv.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, time.Microsecond*500)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

func connectDB(ctx context.Context, log *logs.Logger) *pgx.Conn {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Error(ctx, "Unable to connect to database", "ERROR", err)
		os.Exit(1)
	}

	err = conn.Ping(ctx)
	if err != nil {
		log.Error(ctx, "Connection database", "ERROR", err)
		os.Exit(1)
	}

	return conn
}
