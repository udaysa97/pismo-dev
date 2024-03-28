package email

import (
	"context"
)

type EmailInterface interface {
	SendMail(ctx context.Context, payload map[string]interface{}) (bool, error)
}
