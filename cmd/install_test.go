package cmd

import "testing"

func TestInstallCmdRejectsInvalidIaCToolBinary(t *testing.T) {
	cmd := installCmd()
	cmd.SetArgs([]string{"--iac-tool", "invalid", "hashicorp/aws"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected an error for invalid --iac-tool value")
	}
}
