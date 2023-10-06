package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/peonii/feta/internal/cmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	cmd.Execute(ctx)
}
