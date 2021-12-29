package app

import (
	"go/build"
	"os"
)

func GetCurrentGoPath() string {
	path := os.Getenv("GOPATH")
	if path == "" {
		path = build.Default.GOPATH
	}

	return path
}
