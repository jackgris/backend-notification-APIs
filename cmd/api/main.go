package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"

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

	portEnv := os.Getenv("PORT")

	port, err := strconv.Atoi(portEnv)
	if err != nil {
		log.Error(ctx, "Environment variable PORT converting to integer")
		os.Exit(1)
	}

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0), "Server started at port", port)
	err = http.ListenAndServe(":"+portEnv, mux)
	if err != nil {
		return err
	}

	log.Info(ctx, "Server shutdown OK")

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
