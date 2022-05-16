package web

import (
	"discord-oauth2-example/config"
	"discord-oauth2-example/discord"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

var le *logrus.Entry
var cfg *config.Config

func TestMain(m *testing.M) {
	le = logrus.New().WithField("service", "testing")
	cfg, _ = config.NewConfig()
	os.Exit(m.Run())
}

func TestApi_Redirect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	discordSvc := discord.NewMockService(ctrl)
	discordSvc.EXPECT().AuthCsrfUrl()
	api := NewRestApi(le, cfg, discordSvc)

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	r := httptest.NewRecorder()
	mux := api.GetMux()
	mux.ServeHTTP(r, req)
	assert.Equal(t, r.Result().StatusCode, http.StatusTemporaryRedirect)
}

func TestApi_Auth(t *testing.T) {
	t.Run("missing csrf", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		discordSvc := discord.NewMockService(ctrl)
		api := NewRestApi(le, cfg, discordSvc)
		req := httptest.NewRequest(http.MethodGet, "/redirect/oauth", nil)

		r := httptest.NewRecorder()
		mux := api.GetMux()
		mux.ServeHTTP(r, req)
		assert.Equal(t, r.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("csrf mismatch", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		discordSvc := discord.NewMockService(ctrl)
		api := NewRestApi(le, cfg, discordSvc)
		req := httptest.NewRequest(http.MethodGet, "/redirect/oauth", nil)
		req.AddCookie(&http.Cookie{
			Name:  "csrfToken",
			Value: "1234",
		})
		r := httptest.NewRecorder()
		mux := api.GetMux()
		mux.ServeHTTP(r, req)
		assert.Equal(t, http.StatusBadRequest, r.Result().StatusCode)
	})
	t.Run("bad auth", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		discordSvc := discord.NewMockService(ctrl)
		api := NewRestApi(le, cfg, discordSvc)
		req := httptest.NewRequest(http.MethodGet, "/redirect/oauth", nil)

		discordSvc.EXPECT().Auth(req.Context(), "1234").Return(nil, errors.New(""))
		req.AddCookie(&http.Cookie{
			Name:  "csrfToken",
			Value: "1234",
		})
		reqForm := url.Values{}
		reqForm.Add("state", "1234")
		reqForm.Add("code", "1234")
		req.Form = reqForm
		r := httptest.NewRecorder()
		mux := api.GetMux()
		mux.ServeHTTP(r, req)
		assert.Equal(t, http.StatusInternalServerError, r.Result().StatusCode)
	})
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		discordSvc := discord.NewMockService(ctrl)
		api := NewRestApi(le, cfg, discordSvc)
		req := httptest.NewRequest(http.MethodGet, "/redirect/oauth", nil)
		user := &discord.User{
			UserId:   "1",
			Username: "Linus",
			Email:    "linus@linusjakobsson.com",
		}
		discordSvc.EXPECT().Auth(req.Context(), "1234").Return(user, nil)

		req.AddCookie(&http.Cookie{
			Name:  "csrfToken",
			Value: "1234",
		})
		reqForm := url.Values{}
		reqForm.Add("state", "1234")
		reqForm.Add("code", "1234")
		req.Form = reqForm
		r := httptest.NewRecorder()
		mux := api.GetMux()
		mux.ServeHTTP(r, req)

		decUser := new(discord.User)
		json.NewDecoder(r.Result().Body).Decode(decUser)
		assert.Equal(t, http.StatusOK, r.Result().StatusCode)
		assert.Equal(t, user.UserId, user.UserId)
		assert.Equal(t, decUser.Username, user.Username)
		assert.Equal(t, decUser.Email, user.Email)
	})
}
