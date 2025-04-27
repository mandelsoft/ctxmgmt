package main

import (
	"fmt"
	"os"
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
	default:
		err = fmt.Errorf("unknown example %q", cmd)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
