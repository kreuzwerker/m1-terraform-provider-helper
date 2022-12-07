package cmd

import (
	"github.com/spf13/cobra"
)

func lockfileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lockfile",
		Short: "Commands to work with Terraform lockfiles",
	}

	cmd.AddCommand(lockfileUpgradeCmd())

	return cmd
}
