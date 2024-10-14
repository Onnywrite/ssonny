package main

import (
	"os"
	"os/signal"

	"github.com/Onnywrite/ssonny/internal/app"
	"github.com/Onnywrite/ssonny/internal/config"
	"github.com/Onnywrite/ssonny/pkg/must"
)

func main() {
	cfg := must.Ok2(config.Load("sso.yaml", "/etc/sso/sso.yaml"))
	config.Set(cfg)

	application := app.New()
	application.MustStart()

	shut := make(chan os.Signal, 1)
	signal.Notify(shut, os.Interrupt)
	<-shut

	application.MustStop()
}
