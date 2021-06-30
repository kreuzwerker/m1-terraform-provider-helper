package app_test

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/kreuzwerker/m1-terraform-provider-helper/internal/app"
)

func isDirExistent(dir string) bool {
	_, foldererr := os.Stat(dir)

	return !os.IsNotExist(foldererr)
}

func createDirIfNotExists(dir string) {
	if !isDirExistent(dir) {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func setupTestAppInstance(t *testing.T) (app.App, *bytes.Buffer) {
	tmpDir := t.TempDir()
	buf := new(bytes.Buffer)
	config := app.Config{
		tmpDir + "/plugins", tmpDir + "/plugins_backup", tmpDir, tmpDir + "/cliCache",
	}
	app := app.App{
		Config: &config,
		Out:    buf,
	}
	app.Init()
	return app, buf
}

func TestActivate(t *testing.T) {
	t.Run("Should activate on first use (no directories present)", func(t *testing.T) {
		app, buf := setupTestAppInstance(t)
		app.Activate()
		out := buf.String()
		if out != "Activated\n" {
			t.Fatalf("expected %#v, but got %#v", "Activated\n", out)
		}
	})

	t.Run("Should activate if Plugin Backup directory is present", func(t *testing.T) {
		app, buf := setupTestAppInstance(t)
		createDirIfNotExists(app.Config.TerraformPluginBackupDir)
		app.Activate()
		out := buf.String()
		if out != "Activated\n" {
			t.Fatalf("expected %#v, but got %#v", "Activated\n", out)
		}
	})
	t.Run("Should already be activated if Plugin directory is present", func(t *testing.T) {
		app, buf := setupTestAppInstance(t)
		createDirIfNotExists(app.Config.TerraformPluginDir)
		app.Activate()
		out := buf.String()
		if out != "Already activated\n" {
			t.Fatalf("expected %#v, but got %#v", "Already activated\n", out)
		}
	})
}

func TestDeactivate(t *testing.T) {
	t.Run("Should deactivate on first use (no directories present)", func(t *testing.T) {
		app, buf := setupTestAppInstance(t)
		app.Deactivate()
		out := buf.String()
		if out != "Deactivated\n" {
			t.Fatalf("expected %#v, but got %#v", "Deactivated\n", out)
		}
	})
	t.Run("Should deactivate if Plugin directory is present", func(t *testing.T) {
		app, buf := setupTestAppInstance(t)
		createDirIfNotExists(app.Config.TerraformPluginDir)
		app.Deactivate()
		out := buf.String()
		if out != "Deactivated\n" {
			t.Fatalf("expected %#v, but got %#v", "Deactivated\n", out)
		}
	})
	t.Run("Should already be deactivated if Plugin Backup directory is present", func(t *testing.T) {
		app, buf := setupTestAppInstance(t)
		createDirIfNotExists(app.Config.TerraformPluginBackupDir)
		app.Deactivate()
		out := buf.String()
		if out != "Already Deactivated\n" {
			t.Fatalf("expected %#v, but got %#v", "Already Deactivated\n", out)
		}
	})
}
