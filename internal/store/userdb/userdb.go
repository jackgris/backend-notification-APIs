package userdb

import (
	"context"
	"fmt"
	"time"

	"github.com/jackgris/backend-notification-APIs/internal/domain/usermodel"
	"github.com/jackgris/backend-notification-APIs/internal/store"
)

type Store struct {
	db store.PgxIface
}

func NewStore(db store.PgxIface) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetUsers(category string) ([]usermodel.User, error) {
	var users []usermodel.User

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()

	// Fetch users subscribed to this category
	rows, err := s.db.Query(ctx, `
		SELECT id, name, email, phone, subscribed_categories, notification_channels
		FROM users
		WHERE $1 = ANY(subscribed_categories)
	`, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user usermodel.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.SubscribedCategories, &user.NotificationChannels)
		if err != nil {
			return nil, fmt.Errorf("error row scan: %s", err)
		}
		users = append(users, user)
	}

	return users, nil
}
