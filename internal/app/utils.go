package app

import (
	"go/build"
	"os/exec"
	"strings"
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
