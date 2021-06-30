package cmd

import (
	"github.com/kreuzwerker/m1-terraform-provider-helper/internal/app"
	"github.com/spf13/cobra"
)

func statusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Shows the status of the m1 provider installations",

		RunE: func(cmd *cobra.Command, args []string) error {
			a := app.New()
			a.Init()
			a.CheckStatus()

			return nil
		},
	}

	return cmd
}
