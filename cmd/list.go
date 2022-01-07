package cmd

import (
	"github.com/kreuzwerker/m1-terraform-provider-helper/internal/app"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all available providers and their versions",

		RunE: func(cmd *cobra.Command, args []string) error {
			a := app.New()
			a.Init()
			a.ListProviders()

			return nil
		},
	}

	return cmd
}
