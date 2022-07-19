package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// RootCmd is a root Cobra command that gets called
// from the main func. All other sub-commands should
// be registered here.
func RootCmd() *cobra.Command {
	var v string

	cmd := &cobra.Command{
		Use:   "m1-terraform-provider-helper",
		Short: "A CLI to manage the installation of terraform providers for the Mac M1 chip",
	}

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := setUpLogs(os.Stdout, v); err != nil {
			return err
		}

		return nil
	}

	cmd.PersistentFlags().StringVarP(&v, "verbosity", "v", logrus.WarnLevel.String(), "Log level (debug, info, warn, error, fatal, panic")

	cmd.AddCommand(
		statusCmd(),
		activateCmd(),
		deactivateCmd(),
		lockfileCmd(),
		installCmd(),
		listCmd(),
		versionCmd(),
	)

	return cmd
}

func setUpLogs(out io.Writer, level string) error {
	logrus.SetOutput(out)
	lvl, err := logrus.ParseLevel(level)

	if err != nil {
		return fmt.Errorf("invalid log level: %s", level)
	}

	logrus.SetLevel(lvl)

	return nil
}
