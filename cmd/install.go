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

	var customProviderRepositoryURL string

	var iaCToolBinary string

	cmd := &cobra.Command{
		Use:   "install [providerName]",
		Args:  cobra.ExactArgs(1),
		Short: "Download (and compile) a Terraform/OpenTofu provider for an ARM-based Mac",
		Long:  "Download and compile a specific Terraform/OpenTofu provider for an ARM-based Mac. Provider name is the registry identifier (e.g., \"hashicorp/aws\")",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !isValidIaCToolBinary(iaCToolBinary) {
				return fmt.Errorf("invalid value for --iac-tool: %s", iaCToolBinary)
			}

			a := app.New(iaCToolBinary)
			a.Init()

			if customProviderRepositoryURL != "" {
				a.SetCustomProviderRepositoryURL(customProviderRepositoryURL)
			}

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
	cmd.Flags().StringVar(&iaCToolBinary, "iac-tool", app.DefaultIaCToolBinary, "IaC binary to use: auto, terraform, or tofu")
	cmd.Flags().StringVar(&customBuildCommand, "custom-build-command", "", "A custom build command to execute instead of the built-in commands")
	cmd.Flags().StringVarP(&customTerraformRegistryURL, "custom-terraform-registry-url", "u", "", "A custom URL of the Terraform/OpenTofu registry")
	cmd.Flags().StringVarP(&customProviderRepositoryURL, "custom-provider-repository-url", "p", "", "A custom URL of the provider repository")

	return cmd
}

func isValidIaCToolBinary(value string) bool {
	switch value {
	case app.DefaultIaCToolBinary, app.IaCToolTerraform, app.IaCToolTofu:
		return true
	default:
		return false
	}
}
