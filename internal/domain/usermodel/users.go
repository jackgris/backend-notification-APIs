package usermodel

// User model
type User struct {
	ID                   int
	Name                 string
	Email                string
	Phone                string
	SubscribedCategories []string
	NotificationChannels []string
}
