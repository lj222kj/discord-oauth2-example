package discord

import (
	"context"
	"discord-oauth2-example/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var cfg *config.Config

func TestMain(m *testing.M) {
	cfg, _ = config.NewConfig()
	os.Exit(m.Run())
}

func TestMockService_AuthCsrfUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := New(cfg.DiscordOauthConfig)
	csrf, url := svc.AuthCsrfUrl()
	assert.NotEqual(t, "", csrf)
	assert.NotEqual(t, "", url)
}

func TestMockService_Auth(t *testing.T) {
	t.Run("invalid code", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svc := New(cfg.DiscordOauthConfig)
		_, err := svc.Auth(context.TODO(), "1234")
		assert.NotNil(t, err)
	})
}
