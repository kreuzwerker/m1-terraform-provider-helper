package main

import (
	"log"

	"github.com/kreuzwerker/m1-terraform-provider-helper/cmd"
)

func main() {
	root := cmd.RootCmd()
	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
