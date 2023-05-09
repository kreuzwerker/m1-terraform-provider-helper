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

	var customTerraformRegistryURL string

	cmd := &cobra.Command{
		Use:   "install [providerName]",
		Args:  cobra.ExactArgs(1),
		Short: "Download (and compile) a Terraform provider for an ARM-based Mac",
		Long:  "Download and compile a specific Terraform provider for an ARM-based Mac. Provider name is the Terraform registry identifier (e.g., \"hashicorp/aws\")",
		RunE: func(cmd *cobra.Command, args []string) error {
			a := app.New()
			a.Init()

			if customTerraformRegistryURL != "" {
				a.SetTerraformRegistryURL(customTerraformRegistryURL)
			}

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
	cmd.Flags().StringVarP(&customTerraformRegistryURL, "custom-terraform-registry-url", "u", "", "A custom URL of Terraform registry")

	return cmd
}
