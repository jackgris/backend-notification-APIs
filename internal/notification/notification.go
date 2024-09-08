package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackgris/backend-notification-APIs/internal/loggin"
	"github.com/jackgris/backend-notification-APIs/internal/usermodel"
)

type Store interface {
	GetUsers(category string) ([]usermodel.User, error)
}

type Notification struct {
	db   Store
	logs *loggin.Logs
}

func NewNotification(db Store, logs *loggin.Logs) *Notification {
	return &Notification{
		db:   db,
		logs: logs,
	}
}

func (n *Notification) NotifyUsers(w http.ResponseWriter, r *http.Request) {
	var req NotificationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Message == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if !validateCategories(req.Category) {
		http.Error(w, "Invalid category", http.StatusBadRequest)
		return
	}
	users, err := n.db.GetUsers(req.Category)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	// Log notifications
	for _, user := range users {
		for _, channel := range user.NotificationChannels {
			n.logs.Notification(ctx, user, req.Category, req.Message, channel)
		}
	}

	// Respond to client
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Notifications sent and logged.")
}
