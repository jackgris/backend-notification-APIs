package lognotifier

import "time"

// LogEntry for storing notifications in logs
type LogEntry struct {
	UserID           int       `json:"user_id"`
	Category         string    `json:"category"`
	Message          string    `json:"message"`
	NotificationType string    `json:"notification_type"`
	Timestamp        time.Time `json:"timestamp"`
}
