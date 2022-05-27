package cmd

import "github.com/spf13/cobra"

// RootCmd is a root Cobra command that gets called
// from the main func. All other sub-commands should
// be registered here.
func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "m1-terraform-provider-helper",
		Short: "A CLI to manage the installation of terraform providers for the Mac M1 chip",
	}
	cmd.AddCommand(
		statusCmd(),
		activateCmd(),
		deactivateCmd(),
		installCmd(),
		listCmd(),
		versionCmd(),
	)

	return cmd
}
