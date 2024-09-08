package notification

// NotificationRequest is the request body for notifications
type NotificationRequest struct {
	Category string `json:"category"`
	Message  string `json:"message"`
}
