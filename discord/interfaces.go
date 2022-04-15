package discord

import (
	"context"
)

type Service interface {
	Auth(ctx context.Context, code string) (*User, error)
	AuthCsrfUrl() (string, string)
}
