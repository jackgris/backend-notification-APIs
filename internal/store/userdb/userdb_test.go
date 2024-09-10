package userdb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackgris/backend-notification-APIs/internal/store/userdb"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

// Mock data for testing
var values = [][]any{
	{
		1, "Alice", "alice@example.com", "123-456-7890", []string{"Sports"}, []string{"email"},
	},
	{
		2, "Bob", "bob@example.com", "987-654-3210", []string{"Sports"}, []string{"sms"},
	},
}

func TestGetUsers_Success(t *testing.T) {
	// Create a mock store and mock response
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	defer mock.Close(ctx)

	rows := pgxmock.NewRows([]string{"id", "name", "email", "phone", "subscribed_categories", "notification_channels"}).AddRows(values...)

	mock.ExpectQuery("SELECT").WithArgs("Sports").WillReturnRows(rows)

	store := userdb.NewStore(mock)

	// Call the function and check results
	users, err := store.GetUsers("Sports")

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Equal(t, 2, len(users))
}

func TestGetUsers_NoUsers(t *testing.T) {
	// Create a mock store and mock response
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	defer mock.Close(ctx)

	rows := pgxmock.NewRows([]string{"id", "name", "email", "phone", "subscribed_categories", "notification_channels"}).AddRows([][]any{}...)

	mock.ExpectQuery("SELECT").WithArgs("Sports").WillReturnRows(rows)

	store := userdb.NewStore(mock)

	// Call the function and check results
	users, err := store.GetUsers("Sports")

	assert.NoError(t, err)
	assert.Equal(t, 0, len(users)) // No users should be returned
}

func TestGetUsers_RowError(t *testing.T) {
	// Create a mock store and mock response
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	defer mock.Close(ctx)

	rows := pgxmock.NewRows([]string{"id", "name", "email", "phone", "subscribed_categories", "notification_channels"}).
		AddRows(values...).RowError(0, fmt.Errorf("database error"))

	mock.ExpectQuery("SELECT").WithArgs("Sports").WillReturnRows(rows)

	store := userdb.NewStore(mock)

	// Call the function and check results
	users, err := store.GetUsers("Sports")

	assert.Error(t, err)
	assert.Nil(t, users)
	assert.EqualError(t, err, "error row scan: database error")
}

func TestGetUsers_DatabaseError(t *testing.T) {
	// Create a mock store and mock response
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	defer mock.Close(ctx)

	mock.ExpectQuery("SELECT").WithArgs("Sports").WillReturnError(fmt.Errorf("database error"))

	store := userdb.NewStore(mock)

	// Call the function and check results
	users, err := store.GetUsers("Sports")

	assert.Error(t, err)
	assert.Nil(t, users)
	assert.EqualError(t, err, "database error")
}
