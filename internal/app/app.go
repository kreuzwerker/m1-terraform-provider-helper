package app

import (
	"fmt"
	"log"
	"os"
)

// App implements the core logic of generating passwords.
type App struct {
	terraformProviderDir string
	backupDir            string
	homeDir              string
	cliDir               string
}

const (
	cliDir                      = ".m1-terraform-provider-helper"
	DefaultTerraformProviderDir = "/.terraform.d/plugins"
	DefaultBackupDir            = "/.terraform.d/plugins_backup"
)

func New(terraformProviderDir string, backupDir string) *App {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return &App{
		homeDir:              homeDir,
		terraformProviderDir: homeDir + terraformProviderDir,
		backupDir:            homeDir + backupDir,
		cliDir:               homeDir + "/" + cliDir,
	}
}

func (a *App) isDirExistent(dir string) bool {
	_, foldererr := os.Stat(dir)

	return !os.IsNotExist(foldererr)
}

func (a *App) IsTerraformPluginDirExistent() bool {
	return a.isDirExistent(a.terraformProviderDir)
}

func (a *App) Activate() {
	if a.isDirExistent(a.terraformProviderDir) {
		log.Print("test")
		fmt.Fprintln(os.Stdout, "Already active")
	} else {
		err := os.Rename(a.backupDir, a.terraformProviderDir)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(os.Stdout, "Activated")
	}
}

func (a *App) Deactivate() {
	if a.isDirExistent(a.backupDir) {
		fmt.Fprintln(os.Stdout, "Already Deactivated")
	} else {
		err := os.Rename(a.terraformProviderDir, a.backupDir)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(os.Stdout, "Deactivated")
	}
}
