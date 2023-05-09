package cmd

import (
	"github.com/kreuzwerker/m1-terraform-provider-helper/internal/app"
	"github.com/spf13/cobra"
)

func deactivateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deactivate",
		Short: "Deactivate the m1-terraform-provider-helper",

		RunE: func(cmd *cobra.Command, args []string) error {
			a := app.New()
			a.Init()
			a.Deactivate()

			return nil
		},
	}

	return cmd
}
