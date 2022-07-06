package app

import (
	"log"
	"os"
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
	t.Run("Should create body with one entry", func(t *testing.T) {
		lockfile := &Lockfile{
			Provider: []ProviderConfig{{Name: "test", Version: "1.0.0", Constraints: "1.0.0", Hashes: []string{"h1:test", "h1:test1"}}},
		}
		fileContents := createHclBody(*lockfile)

		stringToCompare := `provider "test" {
  version     = "1.0.0"
  constraints = "1.0.0"
  hashes      = ["h1:test", "h1:test1"]
}
`

		if stringToCompare != fileContents {
			t.Fatalf("expected %#v, but got %#v", stringToCompare, fileContents)
		}
	})
	t.Run("Should create body with two entries", func(t *testing.T) {
		lockfile := &Lockfile{
			Provider: []ProviderConfig{{Name: "test", Version: "1.0.0", Constraints: "1.0.0", Hashes: []string{"h1:test"}}, {Name: "test", Version: "1.0.0", Constraints: "1.0.0", Hashes: []string{"h1:test"}}},
		}
		fileContents := createHclBody(*lockfile)

		stringToCompare := `provider "test" {
  version     = "1.0.0"
  constraints = "1.0.0"
  hashes      = ["h1:test"]
}
provider "test" {
  version     = "1.0.0"
  constraints = "1.0.0"
  hashes      = ["h1:test"]
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
