package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackgris/backend-notification-APIs/internal/loggin"
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

	logs := loggin.NewLogs(db)
	store := userdb.NewStore(db)
	n := notification.NewNotification(store, logs)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /notify", middleware.LogResponse(n.NotifyUsers, log))

	log.Info(ctx, "Server started at :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return err
	}

	return nil
}

func connectDB(ctx context.Context, log *logs.Logger) *pgx.Conn {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Error(ctx, fmt.Sprintf("Unable to connect to database: %v\n", err))
		os.Exit(1)
	}

	err = conn.Ping(ctx)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("Connection err: %v\n", err))
		os.Exit(1)
	}

	return conn
}
