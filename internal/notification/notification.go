package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackgris/backend-notification-APIs/internal/lognotifier"
	"github.com/jackgris/backend-notification-APIs/internal/services/email"
	"github.com/jackgris/backend-notification-APIs/internal/services/pushnotification"
	"github.com/jackgris/backend-notification-APIs/internal/services/sms"
	"github.com/jackgris/backend-notification-APIs/internal/usermodel"
	"github.com/jackgris/backend-notification-APIs/pkg/logs"
)

type Store interface {
	GetUsers(category string) ([]usermodel.User, error)
}

type Notification struct {
	db         Store
	notifyLogs *lognotifier.LogNotifier
	log        *logs.Logger
}

func NewNotification(db Store, notifyLogs *lognotifier.LogNotifier, log *logs.Logger) *Notification {
	return &Notification{
		db:         db,
		notifyLogs: notifyLogs,
		log:        log,
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
	// Send notifications and save logs
	for _, user := range users {
		for _, channel := range user.NotificationChannels {
			switch {
			case channel == "SMS":
				sms.Send(ctx, user.ID, req.Message, n.log)
			case channel == "Email":
				email.Send(ctx, user.ID, req.Message, n.log)
			case channel == "PushNotification":
				pushnotification.Send(ctx, user.ID, req.Message, n.log)
			}
			n.notifyLogs.Notification(ctx, user, req.Category, req.Message, channel)
		}
	}

	// Respond to client
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Notifications sent and logged.")
}
