package cmd

import (
	"fmt"
	"os"

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

			if a.IsTerraformPluginDirExistent() {
				fmt.Fprintln(os.Stdout, "Status: Active")
				fmt.Fprintln(os.Stdout, "Local providers are used")
			} else {
				fmt.Fprintln(os.Stdout, "Status: Not Active")
				fmt.Fprintln(os.Stdout, "All providers are downloaded from the configured registries")
			}

			return nil
		},
	}

	return cmd
}
