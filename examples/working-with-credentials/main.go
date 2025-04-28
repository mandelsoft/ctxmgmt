package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mandelsoft/goutils/general"
)

func main() {
	var err error
	cmd := "application"
	args := os.Args[1:]

	if len(args) > 0 {
		if args[0] == "--dir" {
			args = args[1:]
			if len(args) > 0 {
				os.Chdir(args[0])
				os.Setenv("HOME", general.Must(filepath.Abs(".")))
				args = args[1:]
			}
		}
	}
	if len(args) > 0 {
		cmd = args[0]
	}
	switch cmd {
	case "application":
		err = RunApplication()
	case "basic":
		err = BasicCredentialManagement()
	case "repository":
		err = UsingCredentialsRepositories()
	case "config":
		err = UsingCredentialConfig()
	default:
		err = fmt.Errorf("unknown example %q", cmd)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
