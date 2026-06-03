package app

import (
	"errors"
	"testing"
)

var errNotFound = errors.New("not found")

func TestResolveIaCToolBinary(t *testing.T) {
	t.Run("should prefer terraform when both binaries exist in auto mode", func(t *testing.T) {
		got := resolveIaCToolBinary(DefaultIaCToolBinary, func(string) (string, error) { return "/bin/terraform", nil })
		if got != IaCToolTerraform {
			t.Fatalf("expected %q, got %q", IaCToolTerraform, got)
		}
	})

	t.Run("should fall back to opentofu in auto mode when terraform is missing", func(t *testing.T) {
		lookPath := func(binary string) (string, error) {
			if binary == IaCToolTerraform {
				return "", errNotFound
			}
			return "/bin/tofu", nil
		}
		got := resolveIaCToolBinary(DefaultIaCToolBinary, lookPath)
		if got != IaCToolTofu {
			t.Fatalf("expected %q, got %q", IaCToolTofu, got)
		}
	})

	t.Run("should honor explicit terraform selection", func(t *testing.T) {
		got := resolveIaCToolBinary(IaCToolTerraform, func(string) (string, error) { return "", errNotFound })
		if got != IaCToolTerraform {
			t.Fatalf("expected %q, got %q", IaCToolTerraform, got)
		}
	})

	t.Run("should honor explicit opentofu selection", func(t *testing.T) {
		got := resolveIaCToolBinary(IaCToolTofu, func(string) (string, error) { return "", errNotFound })
		if got != IaCToolTofu {
			t.Fatalf("expected %q, got %q", IaCToolTofu, got)
		}
	})
}

func TestNewUsesSelectedIaCToolBinary(t *testing.T) {
	t.Run("terraform selection sets terraform defaults", func(t *testing.T) {
		app := New(IaCToolTerraform)
		if app.Config.IaCToolBinary != IaCToolTerraform {
			t.Fatalf("expected %q, got %q", IaCToolTerraform, app.Config.IaCToolBinary)
		}
		if app.Config.TerraformPluginDir == "" || app.Config.TerraformRegistryURL != DefaultTerraformRegistryURL {
			t.Fatalf("expected terraform defaults to be configured, got %#v", app.Config)
		}
	})

	t.Run("opentofu selection sets opentofu defaults", func(t *testing.T) {
		app := New(IaCToolTofu)
		if app.Config.IaCToolBinary != IaCToolTofu {
			t.Fatalf("expected %q, got %q", IaCToolTofu, app.Config.IaCToolBinary)
		}
		if app.Config.TerraformPluginDir == "" || app.Config.TerraformRegistryURL != DefaultOpenTofuRegistryURL {
			t.Fatalf("expected opentofu defaults to be configured, got %#v", app.Config)
		}
	})

	t.Run("auto mode resolves once using current environment", func(t *testing.T) {
		original := execLookPath
		defer func() { execLookPath = original }()

		execLookPath = func(binary string) (string, error) {
			if binary == IaCToolTerraform {
				return "", errNotFound
			}
			return "/usr/local/bin/tofu", nil
		}

		app := New()
		if app.Config.IaCToolBinary != IaCToolTofu {
			t.Fatalf("expected %q, got %q", IaCToolTofu, app.Config.IaCToolBinary)
		}
		if app.Config.TerraformRegistryURL != DefaultOpenTofuRegistryURL {
			t.Fatalf("expected opentofu registry url, got %q", app.Config.TerraformRegistryURL)
		}
	})
}

func TestGetIaCToolCommandName(t *testing.T) {
	t.Run("should map terraform mode to terraform command", func(t *testing.T) {
		if got := getIaCToolCommandName(IaCToolTerraform); got != IaCToolTerraform {
			t.Fatalf("expected %q, got %q", IaCToolTerraform, got)
		}
	})

	t.Run("should map opentofu mode to tofu command", func(t *testing.T) {
		if got := getIaCToolCommandName(IaCToolTofu); got != IaCToolTofu {
			t.Fatalf("expected %q, got %q", IaCToolTofu, got)
		}
	})
}
