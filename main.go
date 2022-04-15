package main

import (
	"context"
	"discord-oauth2-example/discord"
	"discord-oauth2-example/global"
	"discord-oauth2-example/web"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	le := logrus.New().WithFields(logrus.Fields{
		"app": "tenant-api",
	})
	cfg, err := global.NewConfig()
	if err != nil {
		le.WithField("error", err.Error()).Fatal("failed to parse env config")
	}

	le.Level = logrus.ErrorLevel
	if cfg.Env == "dev" {
		le.Level = logrus.InfoLevel
	}

	sigChan := make(chan os.Signal, 10)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	le.WithFields(logrus.Fields{
		"app_name": cfg.AppName,
		"env":      cfg.Env,
		"port":     cfg.Port,
	}).Info("application starting")
	discordSvc := discord.New(cfg)
	api := web.NewRestApi(le, cfg, discordSvc)
	go func() {
		if err := api.ListenAndServe(); err != nil {
			le.WithField("error", err.Error()).Errorf("http application shutdown, port %s %v", cfg.Port, err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())

	<-sigChan
	le.Info("received termination, graceful shutdown")
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	api.Shutdown(ctx)
}
