package scaffold

import (
	"context"
	"os"

	"github.com/spf13/cobra"
)

var (
	backendAddress string
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scaffold",
		Short: "scaffold the necessary config to use the GitHub Actions Terraform workflow",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&backendAddress, "backend-url", "https://ffddorf-terraform-backend.fly.dev/", "URL to use as the backend address")

	return cmd
}

func run(ctx context.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := writeBackendConfig(cwd); err != nil {
		return err
	}

	// todo: create github actions workflows

	return nil
}
