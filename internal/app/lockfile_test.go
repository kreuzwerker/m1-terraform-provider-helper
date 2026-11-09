package app

import (
	"log"
	"os"
	"strings"
	"testing"
)

func createFileInDir(dir string, fileName string) {
	file, err := os.Create(dir + "/" + fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}

func TestCreateHclBody(t *testing.T) {
	versionString := "1.0.0"
	t.Run("Should create body with one entry", func(t *testing.T) {
		lockfile := &Lockfile{
			Provider: []ProviderConfig{{Name: "test", Version: "1.0.0", Constraints: &versionString, Hashes: []string{"h1:test", "h1:test1"}}},
		}
		fileContents := createHclBody(*lockfile)

		stringToCompare := `# This file is maintained automatically by "terraform init".
# Manual edits may be lost in future updates.

provider "test" {
  version     = "1.0.0"
  constraints = "1.0.0"
  hashes = [
    "h1:test",
    "h1:test1",
  ]
}
`

		if stringToCompare != fileContents {
			t.Fatalf("expected %#v, but got %#v", stringToCompare, fileContents)
		}
	})
	t.Run("Should create body with two entries", func(t *testing.T) {
		lockfile := &Lockfile{
			Provider: []ProviderConfig{{Name: "test", Version: "1.0.0", Constraints: &versionString, Hashes: []string{"h1:test"}}, {Name: "test", Version: "1.0.0", Constraints: &versionString, Hashes: []string{"h1:test"}}},
		}
		fileContents := createHclBody(*lockfile)

		stringToCompare := `# This file is maintained automatically by "terraform init".
# Manual edits may be lost in future updates.

provider "test" {
  version     = "1.0.0"
  constraints = "1.0.0"
  hashes = [
    "h1:test",
  ]
}
provider "test" {
  version     = "1.0.0"
  constraints = "1.0.0"
  hashes = [
    "h1:test",
  ]
}
`

		if stringToCompare != fileContents {
			t.Fatalf("expected %#v, but got %#v", stringToCompare, fileContents)
		}
	})
	t.Run("Should create body with no constraints", func(t *testing.T) {
		lockfile := &Lockfile{
			Provider: []ProviderConfig{{Name: "test", Version: "1.0.0", Hashes: []string{"h1:test"}}},
		}
		fileContents := createHclBody(*lockfile)

		stringToCompare := `# This file is maintained automatically by "terraform init".
# Manual edits may be lost in future updates.

provider "test" {
  version = "1.0.0"
  hashes = [
    "h1:test",
  ]
}
`

		if stringToCompare != fileContents {
			t.Fatalf("expected %#v, but got %#v", stringToCompare, fileContents)
		}
	})
}

func TestParseOutputLockfilePath(t *testing.T) {
	t.Run("Should return default path if empty", func(t *testing.T) {
		outputLockfilePath := ""
		expectedPath := terraformLockfileName
		actualPath := parseOutputLockfilePath(outputLockfilePath)
		if expectedPath != actualPath {
			t.Fatalf("expected %#v, but got %#v", expectedPath, actualPath)
		}
	})
	t.Run("Should return path if not empty", func(t *testing.T) {
		outputLockfilePath := "test"
		expectedPath := outputLockfilePath
		actualPath := parseOutputLockfilePath(outputLockfilePath)
		if expectedPath != actualPath {
			t.Fatalf("expected %#v, but got %#v", expectedPath, actualPath)
		}
	})
}

func TestGetLockfile(t *testing.T) {
	t.Run("Should return default path if not existent", func(t *testing.T) {
		lockFilePath := "test"
		expectedPath := terraformLockfileName
		actualPath := getLockfile(lockFilePath)
		if expectedPath != actualPath {
			t.Fatalf("expected %#v, but got %#v", expectedPath, actualPath)
		}
	})
	t.Run("Should return path if existent", func(t *testing.T) {
		tmpDir := t.TempDir()
		createFileInDir(tmpDir, "test")
		lockFilePath := tmpDir + "/test"
		expectedPath := tmpDir + "/test"
		actualPath := getLockfile(lockFilePath)
		if expectedPath != actualPath {
			t.Fatalf("expected %#v, but got %#v", expectedPath, actualPath)
		}
	})
}

func TestLockfileProviderStatusForLocalInstall(t *testing.T) {
	config := Lockfile{
		Provider: []ProviderConfig{
			{
				Name:    "registry.terraform.io/hashicorp/template",
				Version: "2.1.1",
				Hashes:  []string{"h1:existing"},
			},
		},
	}

	t.Run("returns missing hash when hash is absent", func(t *testing.T) {
		if lockfileProviderStatusForLocalInstall(config, "registry.terraform.io/hashicorp/template", "2.1.1", "h1:local") != lockfileProviderMissingHash {
			t.Fatal("expected missing hash status")
		}
	})

	t.Run("returns ok when hash is present", func(t *testing.T) {
		if lockfileProviderStatusForLocalInstall(config, "registry.terraform.io/hashicorp/template", "2.1.1", "h1:existing") != lockfileProviderOK {
			t.Fatal("expected existing hash to satisfy lockfile")
		}
	})

	t.Run("returns missing entry when provider is absent", func(t *testing.T) {
		if lockfileProviderStatusForLocalInstall(config, "registry.terraform.io/hashicorp/aws", "5.0.0", "h1:local") != lockfileProviderMissingEntry {
			t.Fatal("expected absent provider status")
		}
	})
}

func TestPrintLockfileUpgradeHintIfNeeded(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = os.Chdir(oldWd)
	})

	pluginDir := tmpDir + "/plugins"
	providerAddress := "registry.terraform.io/hashicorp/template"
	version := "2.1.1"
	providerDir := providerInstallDir(pluginDir, providerAddress, version)

	if err := os.MkdirAll(providerDir, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(providerDir+"/terraform-provider-template_v2.1.1_x5", []byte("provider-binary"), 0o644); err != nil {
		t.Fatal(err)
	}

	localHash := getCalculatedHashForProvider(providerDir)

	lockfileContents := `# This file is maintained automatically by "terraform init".
# Manual edits may be lost in future updates.

provider "` + providerAddress + `" {
  version = "` + version + `"
  hashes = [
    "h1:other",
  ]
}
`
	if err := os.WriteFile(terraformLockfileName, []byte(lockfileContents), 0o644); err != nil {
		t.Fatal(err)
	}

	var output strings.Builder
	a := &App{
		Config: &Config{
			TerraformPluginDir:   pluginDir,
			TerraformRegistryURL: DefaultTerraformRegistryURL,
		},
		Out: &output,
	}

	a.printLockfileUpgradeHintIfNeeded("hashicorp/template", version)

	hint := output.String()
	if !strings.Contains(hint, "lockfile upgrade") {
		t.Fatalf("expected lockfile upgrade hint, got: %q", hint)
	}
	if strings.Contains(hint, "terraform providers lock") {
		t.Fatalf("did not expect terraform providers lock in hint, got: %q", hint)
	}

	lockfileContents = `# This file is maintained automatically by "terraform init".
# Manual edits may be lost in future updates.

provider "` + providerAddress + `" {
  version = "` + version + `"
  hashes = [
    "` + localHash + `",
  ]
}
`
	if err := os.WriteFile(terraformLockfileName, []byte(lockfileContents), 0o644); err != nil {
		t.Fatal(err)
	}

	output.Reset()
	a.printLockfileUpgradeHintIfNeeded("hashicorp/template", version)

	if output.Len() != 0 {
		t.Fatalf("expected no hint when hash is present, got: %q", output.String())
	}

	lockfileContents = `not valid hcl for terraform lock file
`
	if err := os.WriteFile(terraformLockfileName, []byte(lockfileContents), 0o644); err != nil {
		t.Fatal(err)
	}

	output.Reset()
	a.printLockfileUpgradeHintIfNeeded("hashicorp/template", version)

	if output.Len() != 0 {
		t.Fatalf("expected no hint on lockfile parse error, got: %q", output.String())
	}

	if err := os.Remove(terraformLockfileName); err != nil {
		t.Fatal(err)
	}

	lockfileContents = `# This file is maintained automatically by "terraform init".

provider "registry.terraform.io/hashicorp/aws" {
  version = "5.0.0"
  hashes = [
    "h1:other",
  ]
}
`
	if err := os.WriteFile(terraformLockfileName, []byte(lockfileContents), 0o644); err != nil {
		t.Fatal(err)
	}

	output.Reset()
	a.printLockfileUpgradeHintIfNeeded("hashicorp/template", version)

	hint = output.String()
	if !strings.Contains(hint, "has no entry for locally installed provider") {
		t.Fatalf("expected missing provider entry hint, got: %q", hint)
	}
	if strings.Contains(hint, "does not include a checksum") {
		t.Fatalf("expected provider-entry hint, not checksum hint, got: %q", hint)
	}
}

func TestParseHclLockfile(t *testing.T) {
	t.Run("Should parse lockfile with no constraints field", func(t *testing.T) {
		parseHclLockfile("../../testdata/lockfile_no_constraints.hcl")
	})
	t.Run("Should parse lockfile with all fields", func(t *testing.T) {
		parseHclLockfile("../../testdata/lockfile.hcl")
	})
}
