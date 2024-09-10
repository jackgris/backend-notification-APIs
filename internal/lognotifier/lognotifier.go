package lognotifier

import (
	"context"
	"log"
	"time"

	"github.com/jackgris/backend-notification-APIs/internal/domain/usermodel"
	"github.com/jackgris/backend-notification-APIs/internal/store"
)

type LogNotifier struct {
	db store.PgxIface
}

func NewLogs(db store.PgxIface) *LogNotifier {
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
