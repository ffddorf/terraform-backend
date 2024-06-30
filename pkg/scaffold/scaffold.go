package scaffold

import (
	"context"
	"io"

	"github.com/nimbolus/terraform-backend/pkg/fs"
	"github.com/spf13/cobra"
)

var (
	backendAddress string
)

func NewCommand(dir fs.FS, stdin io.Reader) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scaffold",
		Short: "scaffold the necessary config to use the GitHub Actions Terraform workflow",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), dir, stdin)
		},
	}

	cmd.Flags().StringVar(&backendAddress, "backend-url", "https://ffddorf-terraform-backend.fly.dev/", "URL to use as the backend address")

	return cmd
}

func run(ctx context.Context, dir fs.FS, stdin io.Reader) error {
	if err := writeBackendConfig(ctx, dir, stdin); err != nil {
		return err
	}

	if err := writeGithubActionsWorkflows(ctx, dir, stdin); err != nil {
		return err
	}

	return nil
}
