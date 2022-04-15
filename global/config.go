package global

import (
	"context"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	AppName           string `env:"APP_NAME,default=discord-oauth-example"`
	Env               string `env:"ENV,default=dev"`
	Host              string `env:"HOST,default=127.0.0.1"`
	Port              string `env:"PORT,default=8080"`
	OAuthRedirectUrl  string `env:"OAUTH_REDIRECT_URL,default=http://localhost:8080/redirect/oauth"`
	OAuthClientSecret string `env:"OAUTH_CLIENT_SECRET,required"`
	OAuthClientId     string `env:"OAUTH_CLIENT_ID,required"`
	AuthUrl           string `env:"DISCORD_AUTH_URL,default=https://discord.com/api/oauth2/authorize"`
	TokenUrl          string `env:"DISCORD_TOKEN_URL,default=https://discord.com/api/oauth2/token"`
	RevokeTokenUrl    string `env:"DISCORD_REVOKE_URL,default=https://discord.com/api/oauth2/token/revoke"`
}

func NewConfig() (*Config, error) {
	ctx := context.Background()
	config := new(Config)
	if err := envconfig.Process(ctx, config); err != nil {
		return nil, err
	}

	return config, nil
}
