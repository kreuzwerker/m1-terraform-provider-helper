package cmd

import (
	"github.com/kreuzwerker/m1-terraform-provider-helper/internal/app"
	"github.com/spf13/cobra"
)

func lockfileUpgradeCmd() *cobra.Command {
	var inputLockfilePath string

	var outputPath string

	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrades the hashes in a terraform lockfile",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {

			a := app.New()
			a.Init()
			a.UpgradeLockfile(inputLockfilePath, outputPath)

			return nil
		},
	}
	cmd.Flags().StringVar(&inputLockfilePath, "input-lockfile-path", "", "Path of lockfile to upgrade. Defaults to .terraform.lock.hcl in the current directory")
	cmd.Flags().StringVar(&outputPath, "output-path", "", "Path of output file")

	return cmd
}
