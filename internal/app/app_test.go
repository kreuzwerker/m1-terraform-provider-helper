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

func TestIsTerraformPluginDirExistent(t *testing.T) {
	t.Run("Should return true if TerraformPluginDir exists", func(t *testing.T) {
		app, _ := setupTestAppInstance(t)
		createDirIfNotExists(app.Config.TerraformPluginDir)
		exists := app.IsTerraformPluginDirExistent()
		if !exists {
			t.Fatalf("Expected %s to exist", app.Config.TerraformPluginDir)
		}
	})
	t.Run("Should return false if TerraformPluginDir exists", func(t *testing.T) {
		app, _ := setupTestAppInstance(t)
		exists := app.IsTerraformPluginDirExistent()
		if exists {
			t.Fatalf("Expected %s to exist", app.Config.TerraformPluginDir)
		}
	})
}

func TestNew(t *testing.T) {
	t.Run("Should create App instance with correct paths", func(t *testing.T) {
		app := app.New()
		baseDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		if app.Config.BaseDir != baseDir {
			t.Fatalf("expected %v, but got %v", baseDir, app.Config.BaseDir)
		}
	})
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

func TestCheckStatus(t *testing.T) {
	t.Run("Should return active status", func(t *testing.T) {
		app, buf := setupTestAppInstance(t)
		createDirIfNotExists(app.Config.TerraformPluginDir)
		app.CheckStatus()
		out := buf.String()
		if out != "Status: Active\nLocal providers are used\n" {
			t.Fatalf("expected %#v, but got %#v", "Status: Active\nLocal providers are used\n", out)
		}
	})
	t.Run("Should return Deactivated status", func(t *testing.T) {
		app, buf := setupTestAppInstance(t)
		app.CheckStatus()
		out := buf.String()
		if out != "Status: Not Active\nAll providers are downloaded from the configured registries\n" {
			t.Fatalf("expected %#v, but got %#v", "Status: Not Active\nAll providers are downloaded from the configured registries\n", out)
		}
	})
}
