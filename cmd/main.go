package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/Onnywrite/ssonny/internal/app"
	"github.com/Onnywrite/ssonny/pkg/must"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	must.Ok1(app.New().Run(ctx))
}
