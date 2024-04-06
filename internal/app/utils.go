package app

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"strings"

	"github.com/adrg/xdg"
)

func GetCurrentGoPath() string {
	buff := &strings.Builder{}
	cmd := exec.Command("go", "env", "GOPATH")
	cmd.Stdout = buff

	if err := cmd.Run(); err != nil {
		panic(err)
	}

	path := strings.TrimSpace(buff.String())
	if path == "" {
		path = build.Default.GOPATH
	}

	return path
}

func GetTerraformPluginsPath(baseDir string) string {
	return getTerraformDataPath(baseDir) + "/plugins"
}

func GetTerraformPluginsBackupPath(baseDir string) string {
	return getTerraformDataPath(baseDir) + "/plugins_backup"
}

func GetProvidersCachePath(baseDir string) string {
	relPath := "m1-terraform-provider-helper"
	cacheDirPath, err := xdg.SearchCacheFile(relPath)
	if err != nil {
		return fmt.Sprintf("%s/.%s", baseDir, relPath)
	}
	return cacheDirPath
}

func GetTerraformRegistryURL() string {
	registryURL, ok := os.LookupEnv("TF_HELPER_REGISTRY_URL")
	if !ok {
		return DefaultTerraformRegistryURL
	}
	return registryURL
}

func getTerraformDataPath(baseDir string) string {
	dataDirPath, err := xdg.SearchDataFile("terraform")
	if err != nil {
		dataDirPath = baseDir + "/.terraform.d"
	}
	return dataDirPath
}
