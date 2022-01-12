package app

import (
	"testing"
)

func TestExtractVersionAsNumber(t *testing.T) {
	var number int
	number = extractMajorVersionAsNumber("v2.34.5")
	if number != 2 {
		t.Fatalf("number should be equal 2")
	}
	number = extractMajorVersionAsNumber("2.34.5")
	if number != 2 {
		t.Fatalf("number should be equal 2")
	}
}

func TestCreateBuildCommand(t *testing.T) {
	var buildCommand string
	var expectedBuildCommand string
	buildCommand = createBuildCommand("datadog/datadog", "2.35.9", "")
	expectedBuildCommand = "make build"
	if buildCommand != expectedBuildCommand {
		t.Fatalf("buildCommand '%s' should be equal '%s'", buildCommand, expectedBuildCommand)
	}
	buildCommand = createBuildCommand("hashicorp/aws", "v2.71.0", "")
	expectedBuildCommand = "make tools && make fmt && gofmt -s -w ./tools.go && make build"
	if buildCommand != expectedBuildCommand {
		t.Fatalf("buildCommand '%s' should be equal '%s'", buildCommand, expectedBuildCommand)
	}
	buildCommand = createBuildCommand("hashicorp/aws", "v3.0.0", "")
	expectedBuildCommand = "cd tools && go get -d github.com/pavius/impi/cmd/impi && cd .. && make tools && make build"
	if buildCommand != expectedBuildCommand {
		t.Fatalf("buildCommand '%s' should be equal '%s'", buildCommand, expectedBuildCommand)
	}
}

func TestNormalizeSemver(t *testing.T) {
	version := normalizeSemver("v2.34.5")
	if version != "2.34.5" {
		t.Fatalf("version should be equal 2.34.5, not %s", version)
	}
	version = normalizeSemver("2.34.5")
	if version != "2.34.5" {
		t.Fatalf("version2 should be equal 2.34.5, not %s", version)
	}
}
