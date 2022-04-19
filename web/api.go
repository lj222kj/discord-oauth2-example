package web

import (
	"context"
	"discord-oauth2-example/config"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Api struct {
	discordSvc DiscordService
	le         *logrus.Entry
	cfg        *config.Config
	srv        *http.Server
	mux        *http.ServeMux
}

func (a *Api) Redirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		csrfToken, url := a.discordSvc.AuthCsrfUrl()
		c := &http.Cookie{
			Name:     "csrfToken",
			Value:    csrfToken,
			Expires:  time.Time{},
			MaxAge:   15,
			Secure:   true,
			HttpOnly: true,
			SameSite: 0,
			Raw:      "",
			Unparsed: nil,
		}
		http.SetCookie(w, c)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)

	}
}

func (a *Api) Auth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("csrfToken")
		fmt.Println(c)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.FormValue("state") != c.Value {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("State does not match."))
			return
		}
		u, err := a.discordSvc.Auth(r.Context(), r.FormValue("code"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ub, err := json.Marshal(u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(ub)
	}
}

func (a *Api) GetMux() *http.ServeMux {
	return a.mux
}
func NewRestApi(le *logrus.Entry, cfg *config.Config, discordSvc DiscordService) *Api {
	mux := http.NewServeMux()

	api := &Api{
		discordSvc: discordSvc,
		le:         le,
		cfg:        cfg,
		mux:        mux,
		srv: &http.Server{
			Addr:         cfg.Host + ":" + cfg.Port,
			Handler:      mux,
			ReadTimeout:  time.Second * 15,
			WriteTimeout: time.Second * 15,
			IdleTimeout:  time.Second * 30,
		},
	}

	mux.Handle("/", api.Redirect())
	mux.Handle("/redirect/oauth", api.Auth())

	return api
}

func (a Api) ListenAndServe() error {
	return a.srv.ListenAndServe()
}

func (a Api) Shutdown(ctx context.Context) error {
	return a.srv.Shutdown(ctx)
}
