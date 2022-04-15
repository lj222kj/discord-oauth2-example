package discord

import (
	"context"
	"discord-oauth2-example/global"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var cfg *global.Config

func TestMain(m *testing.M) {
	cfg, _ = global.NewConfig()
	os.Exit(m.Run())
}

func TestMockService_AuthCsrfUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := New(cfg)
	csrf, url := svc.AuthCsrfUrl()
	assert.NotEqual(t, "", csrf)
	assert.NotEqual(t, "", url)
}

func TestMockService_Auth(t *testing.T) {
	t.Run("invalid code", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svc := New(cfg)
		_, err := svc.Auth(context.TODO(), "1234")
		assert.NotNil(t, err)
	})
}
