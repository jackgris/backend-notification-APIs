package email

import (
	"context"

	"github.com/jackgris/backend-notification-APIs/pkg/logs"
)

func Send(ctx context.Context, id int, message string, log *logs.Logger) {
	// TODO: Should implement the service itself, with retry logic and an error-handling policy.
	log.Info(ctx, "EMAIL", "USER", id, "message", message)
}
