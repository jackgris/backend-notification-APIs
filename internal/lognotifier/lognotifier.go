package lognotifier

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackgris/backend-notification-APIs/internal/usermodel"
)

type LogNotifier struct {
	db *pgx.Conn
}

func NewLogs(db *pgx.Conn) *LogNotifier {
	return &LogNotifier{
		db: db,
	}
}

// Log notifications to the database
func (l *LogNotifier) Notification(ctx context.Context, user usermodel.User, category string, message string, channel string) {
	_, err := l.db.Exec(ctx, `INSERT INTO logs (user_id, category, message, notification_type, timestamp)
		VALUES ($1, $2, $3, $4, $5)`,
		user.ID, category, message, channel, time.Now())

	if err != nil {
		log.Println("Failed to log notification:", err)
	}
}
