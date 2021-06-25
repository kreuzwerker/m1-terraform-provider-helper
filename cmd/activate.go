package cmd

import (
	"github.com/kreuzwerker/m1-terraform-provider-helper/internal/app"
	"github.com/spf13/cobra"
)

func activateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activate",
		Short: "Activate the usage of m1 provider helper",

		RunE: func(cmd *cobra.Command, args []string) error {
			a := app.New(app.DefaultTerraformProviderDir, app.DefaultBackupDir)
			a.Activate()

			return nil
		},
	}

	return cmd
}
