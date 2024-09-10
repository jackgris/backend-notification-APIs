package notification_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pashagolub/pgxmock/v4"

	"github.com/jackgris/backend-notification-APIs/internal/domain/usermodel"
	"github.com/jackgris/backend-notification-APIs/internal/lognotifier"
	"github.com/jackgris/backend-notification-APIs/internal/notification"
	"github.com/jackgris/backend-notification-APIs/internal/store/userdb"
	"github.com/jackgris/backend-notification-APIs/pkg/logs"
)

type testCase struct {
	name         string
	reqBody      string
	statusCode   int
	responseBody string
	hasError     bool             // Flag to indicate expected error in validation or database call
	mockUsers    []usermodel.User // Mock users for database retrieval
}

// Mock data for testing
var values = [][]any{
	{
		1, "Alice", "alice@example.com", "123-456-7890", []string{"Sports"}, []string{"email"},
	},
	{
		2, "Bob", "bob@example.com", "987-654-3210", []string{"Sports"}, []string{"sms"},
	},
}

func TestNotifyUsers(t *testing.T) {
	var testCases []testCase

	// Valid request
	testCases = append(testCases, testCase{
		name:         "Valid Request",
		reqBody:      `{"message": "Test message", "category": "Finance"}`,
		statusCode:   http.StatusOK,
		responseBody: "Notifications sent and logged.",
		mockUsers:    []usermodel.User{{ID: 1, NotificationChannels: []string{"SMS"}}},
	})

	// Empty message
	testCases = append(testCases, testCase{
		name:         "Empty message",
		reqBody:      `{"category": "Films"}`,
		statusCode:   http.StatusBadRequest,
		responseBody: "Invalid request",
	})

	// Invalid category
	testCases = append(testCases, testCase{
		name:         "Invalid category",
		reqBody:      `{"message": "Test message", "category": "invalid_category"}`,
		statusCode:   http.StatusBadRequest,
		responseBody: "Invalid category",
		mockUsers:    []usermodel.User{{ID: 1, NotificationChannels: []string{"SMS"}}},
	})

	// Database error
	testCases = append(testCases, testCase{
		name:         "Database Error",
		reqBody:      `{"message": "Test message", "category": "Sports"}`,
		statusCode:   http.StatusInternalServerError,
		responseBody: http.StatusText(http.StatusInternalServerError),
		hasError:     true,
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// Mock database retrieval based on test case flag
			// Create a mock store and mock response
			mock, err := pgxmock.NewConn()
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.Background()

			defer mock.Close(ctx)
			if !tc.hasError {
				rows := pgxmock.NewRows([]string{"id", "name", "email", "phone", "subscribed_categories", "notification_channels"}).AddRows(values...)
				mock.ExpectQuery("SELECT").WithArgs("Finance").WillReturnRows(rows)
			}

			log := logs.New(io.Discard)

			lognotifier := lognotifier.NewLogs(mock)
			store := userdb.NewStore(mock)
			n := notification.NewNotification(store, lognotifier, log)

			// Create request and recorder
			req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(tc.reqBody))
			if err != nil {
				t.Fatal(err)
			}
			rec := httptest.NewRecorder()

			// Call the function
			n.NotifyUsers(rec, req)

			// Check status code
			if rec.Code != tc.statusCode {
				t.Errorf("Expected status code %d, got %d", tc.statusCode, rec.Code)
			}

			// Check body response
			if body := strings.TrimSpace(rec.Body.String()); body != tc.responseBody {
				t.Errorf("Expected response body %s, got %s", tc.responseBody, body)
			}
		})
	}
}
