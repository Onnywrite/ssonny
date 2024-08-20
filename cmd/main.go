package main

import (
	"os"
	"os/signal"

	"github.com/Onnywrite/ssonny/internal/app"
	"github.com/Onnywrite/ssonny/internal/config"
)

func main() {
	cfg := config.MustLoad("/config/ignore-config.yaml")

	application := app.New(cfg)
	application.MustStart()

	shut := make(chan os.Signal, 1)
	signal.Notify(shut, os.Interrupt)
	<-shut

	application.MustStop()
}
