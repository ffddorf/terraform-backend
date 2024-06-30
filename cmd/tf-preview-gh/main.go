package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nimbolus/terraform-backend/pkg/fs"
	"github.com/nimbolus/terraform-backend/pkg/scaffold"
	"github.com/nimbolus/terraform-backend/pkg/speculative"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("failed to get working directory: %w", err))
	}

	rootCmd := speculative.NewCommand()
	rootCmd.AddCommand(scaffold.NewCommand(fs.ForOS(cwd), os.Stdin))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
