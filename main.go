package main

import (
	"context"
	"discord-oauth2-example/config"
	"discord-oauth2-example/discord"
	"discord-oauth2-example/web"
	"github.com/sirupsen/logrus"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	le := logrus.New().WithFields(logrus.Fields{
		"app": "discord-oauth2-example",
	})

	cfg, err := config.NewConfig()
	if err != nil {
		le.WithError(err).Fatal("failed to parse env config")
	}

	le.Level = logrus.ErrorLevel
	if cfg.Env == "dev" {
		le.Level = logrus.InfoLevel
	}

	le.WithFields(logrus.Fields{
		"app_name": cfg.AppName,
		"env":      cfg.Env,
		"port":     cfg.Port,
	}).Info("application starting")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	discordSvc := discord.New(cfg.DiscordOauthConfig)
	api := web.NewRestApi(le, cfg, discordSvc)

	go func() {
		err := api.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			le.WithError(err).Errorf("http application shutdown, port %s %v\"", cfg.Port, err)
		}
	}()

	<-ctx.Done()
	stop()

	le.Info("graceful shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := api.Shutdown(ctx); err != nil {
		le.Fatalf("Server forced to shutdown: %e", err)
	}

}
