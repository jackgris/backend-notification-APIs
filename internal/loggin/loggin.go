package loggin

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackgris/backend-notification-APIs/internal/usermodel"
)

type Logs struct {
	db *pgx.Conn
}

func NewLogs(db *pgx.Conn) *Logs {
	return &Logs{
		db: db,
	}
}

// Log notifications to the database
func (l *Logs) Notification(ctx context.Context, user usermodel.User, category string, message string, channel string) {
	_, err := l.db.Exec(ctx, `INSERT INTO logs (user_id, category, message, notification_type, timestamp)
		VALUES ($1, $2, $3, $4, $5)`,
		user.ID, category, message, channel, time.Now())

	if err != nil {
		log.Println("Failed to log notification:", err)
	}
}
