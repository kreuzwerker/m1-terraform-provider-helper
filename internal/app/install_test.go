package app

import (
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

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

func TestExtractRepoNameFromUrl(t *testing.T) {
	repoDir := extractRepoNameFromURL("https://github.com/hashicorp/terraform-provider-github")
	if repoDir != "terraform-provider-github" {
		t.Fatalf("repoDir should be equal terraform-provider-github, not %s", repoDir)
	}
	repoDir = extractRepoNameFromURL("git@github.com:hashicorp/terraform-provider-github")
	if repoDir != "terraform-provider-github" {
		t.Fatalf("repoDir should be equal terraform-provider-github, not %s", repoDir)
	}
}

func TestCloneRepo(t *testing.T) {
	tmpDir := t.TempDir()
	fullPath := tmpDir + "/terraform-provider-random"
	cloneRepo("https://github.com/hashicorp/terraform-provider-random", fullPath)
	if !isDirExistent(fullPath) {
		t.Fatalf("terraform-provider-random should be a dir inside %s", tmpDir)
	}
}

func TestGetProviderData(t *testing.T) {
	repo := "https://github.com/hashicorp/terraform-provider-aws"
	description := "terraform-provider-aws"

	t.Run("Should get and parse JSON response", func(t *testing.T) {
		provider := "hashicorp/aws"

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", "https://registry.terraform.io/v1/providers/"+provider,
			httpmock.NewStringResponder(200, `{"description": "`+description+`",
			"source": "`+repo+`"}`))
		providerData, err := getProviderData(provider, 2, "")
		if err != nil {
			t.Errorf("Should not have an error %s ", err)
		}
		if providerData.Repo != repo {
			t.Errorf("expected %#v, but got %#v", repo, providerData.Repo)
		}
		if providerData.Description != description {
			t.Errorf("expected %#v, but got %#v", description, providerData.Description)
		}
	})
	t.Run("Should get request timeout error", func(t *testing.T) {
		provider := "hashicorp/google"
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", "https://registry.terraform.io/v1/providers/"+provider,
			httpmock.NewStringResponder(200, `{"description": "`+description+`",
			"source": "`+repo+`"}`).Delay(3*time.Second))
		_, err := getProviderData(provider, 2, "")
		if !strings.HasPrefix(err.Error(), "timeout error") {
			t.Errorf("Expected \"error timeout\" but got %#v", err.Error())
		}
		if err == nil {
			t.Error("Should have an error")
		}
	})
	t.Run("Should error with mismatched JSON", func(t *testing.T) {
		provider := "hashicorp/vault"
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", "https://registry.terraform.io/v1/providers/"+provider,
			httpmock.NewStringResponder(200, `{"test:"812}`))
		_, err := getProviderData(provider, 2, "")
		if err == nil {
			t.Error("Should run into JSON parse error")
		}
	})
}

func TestCheckoutSourceCode(t *testing.T) {
	t.Run("Should checkout source code once", func(t *testing.T) {
		tmpDir := t.TempDir()
		checkoutSourceCode(tmpDir+"/cliCache", "https://github.com/hashicorp/terraform-provider-random", "v2.2.0")
	})
	t.Run("Should checkout two versions of same source code once", func(t *testing.T) {
		tmpDir := t.TempDir()
		checkoutSourceCode(tmpDir+"/cliCache", "https://github.com/hashicorp/terraform-provider-random", "v2.3.1")
		checkoutSourceCode(tmpDir+"/cliCache", "https://github.com/hashicorp/terraform-provider-random", "v2.2.0")
	})
}

func TestParseBuildOutputAndGetBinaryOutputPath(t *testing.T) {
	t.Run("Should parse correct output from go build -o command", func(t *testing.T) {
		expected := "build/darwin_arm64/terraform-provider-pingdom_v1.1.3"
		actual, ok := parseBuildOutputAndGetBinaryOutputPath("go build -o build/darwin_arm64/terraform-provider-pingdom_v1.1.3 ")
		if ok && expected != actual {
			t.Fatalf("expected %#v, but got %#v", expected, actual)
		}
	})
	t.Run("Should parse correct output from go install command", func(t *testing.T) {
		_, ok := parseBuildOutputAndGetBinaryOutputPath("go install")
		if ok {
			t.Fatalf("expected ok to be false")
		}
	})
}
