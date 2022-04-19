package discord

import (
	"context"
	"discord-oauth2-example/config"
	"discord-oauth2-example/utils"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type discord struct {
	scopes   []string
	oauthCfg *oauth2.Config
}

func New(cfg *config.DiscordOauthConfig) *discord {
	return &discord{
		oauthCfg: &oauth2.Config{
			ClientID:     cfg.OAuthClientId,
			ClientSecret: cfg.OAuthClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:   cfg.AuthUrl,
				TokenURL:  cfg.TokenUrl,
				AuthStyle: oauth2.AuthStyleInParams,
			},
			RedirectURL: cfg.OAuthRedirectUrl,
			Scopes: []string{
				"identity",
			},
		},
	}
}

func (d *discord) AuthCsrfUrl() (string, string) {
	csrfToken := uuid.New().String()
	return csrfToken, d.oauthCfg.AuthCodeURL(csrfToken)
}

func (a *discord) Auth(ctx context.Context, code string) (*User, error) {
	t, err := a.oauthCfg.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	if !t.Valid() {
		return nil, fmt.Errorf("invalid token %e", err)
	}
	res, err := a.oauthCfg.Client(ctx, t).Get("https://discord.com/api/users/@me")
	if err != nil {
		return nil, err
	}

	user := new(User)

	if err := utils.ParseBodyResponse(res, user); err != nil {
		return nil, err
	}

	return &User{
		UserId:   user.UserId,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
