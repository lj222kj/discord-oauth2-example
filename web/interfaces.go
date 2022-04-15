package web

import (
	"context"
	"discord-oauth2-example/discord"
)

type DiscordService interface {
	Auth(ctx context.Context, code string) (*discord.User, error)
	AuthCsrfUrl() (string, string)
}
