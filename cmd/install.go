package cmd

import (
	"fmt"
	"os"

	"github.com/kreuzwerker/m1-terraform-provider-helper/internal/app"
	"github.com/spf13/cobra"
)

func installCmd() *cobra.Command {
	var versionString string

	var customBuildCommand string

	cmd := &cobra.Command{
		Use:   "install [providerName]",
		Args:  cobra.ExactArgs(1),
		Short: "Downloads (and compiles) a terraform provider for the M1 chip",
		Long:  "Download and compiles the specifiec terraform provider for your M1 chip. Provider name is the terraform registry identifier, e.g. \"hashicorp/aws\"",
		RunE: func(cmd *cobra.Command, args []string) error {
			a := app.New()
			a.Init()

			if a.IsTerraformPluginDirExistent() {
				a.Install(args[0], versionString, customBuildCommand)
			} else {
				fmt.Fprintln(os.Stdout, "Please activate first")
			}

			return nil
		},
	}
	cmd.Flags().StringVarP(&versionString, "version", "v", "", "The version of the provider")
	cmd.Flags().StringVar(&customBuildCommand, "custom-build-command", "", "A custom build command to execute instead of the built-in commands")

	return cmd
}
