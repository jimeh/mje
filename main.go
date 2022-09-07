package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jimeh/mj2n/commands"
)

func main() {
	cmd, err := commands.NewMJ2N()
	if err != nil {
		fatal(err)
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	err = cmd.ExecuteContext(ctx)
	if err != nil {
		defer os.Exit(1)

		return
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	os.Exit(1)
}
