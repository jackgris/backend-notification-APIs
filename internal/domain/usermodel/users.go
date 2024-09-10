package usermodel

// User model
type User struct {
	ID                   int      `json:"id"`
	Name                 string   `json:"name"`
	Email                string   `json:"email"`
	Phone                string   `json:"phone"`
	SubscribedCategories []string `json:"subscribed"`
	NotificationChannels []string `json:"channels"`
}
