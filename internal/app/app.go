package app

import (
	"fmt"
	"io"
	"log"
	"os"
)

type App struct {
	Config *Config
	Out    io.Writer
}

type Config struct {
	TerraformPluginDir       string
	TerraformPluginBackupDir string
	BaseDir                  string
	ProvidersCacheDir        string
}

const (
	DefaultProvidersCacheDir        = "/.m1-terraform-provider-helper"
	DefaultTerraformPluginDir       = "/.terraform.d/plugins"
	DefaultTerraformPluginBackupDir = "/.terraform.d/plugins_backup"
)

func New() *App {
	BaseDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	app := &App{
		Config: &Config{
			BaseDir:                  BaseDir,
			TerraformPluginDir:       BaseDir + DefaultTerraformPluginDir,
			TerraformPluginBackupDir: BaseDir + DefaultTerraformPluginBackupDir,
			ProvidersCacheDir:        BaseDir + DefaultProvidersCacheDir,
		},
		Out: os.Stdout,
	}

	return app
}

func (a *App) Init() {
	a.createDirIfNotExists(a.Config.ProvidersCacheDir)
}

func (a *App) createDirIfNotExists(dir string) {
	if !a.isDirExistent(dir) {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (a *App) isDirExistent(dir string) bool {
	_, foldererr := os.Stat(dir)

	return !os.IsNotExist(foldererr)
}

func (a *App) IsTerraformPluginDirExistent() bool {
	return a.isDirExistent(a.Config.TerraformPluginDir)
}

func (a *App) Activate() {
	if a.isDirExistent(a.Config.TerraformPluginDir) {
		fmt.Fprintln(a.Out, "Already activated")
	} else {
		if !a.isDirExistent(a.Config.TerraformPluginBackupDir) {
			a.createDirIfNotExists(a.Config.TerraformPluginBackupDir)
		}
		err := os.Rename(a.Config.TerraformPluginBackupDir, a.Config.TerraformPluginDir)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(a.Out, "Activated")
	}
}

func (a *App) Deactivate() {
	if a.isDirExistent(a.Config.TerraformPluginBackupDir) {
		fmt.Fprintln(a.Out, "Already Deactivated")
	} else {
		if !a.isDirExistent(a.Config.TerraformPluginDir) {
			a.createDirIfNotExists(a.Config.TerraformPluginDir)
		}
		err := os.Rename(a.Config.TerraformPluginDir, a.Config.TerraformPluginBackupDir)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(a.Out, "Deactivated")
	}
}
