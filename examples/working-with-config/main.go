package main

import (
	"fmt"
	"os"
)

var current_version string

func main() {
	var err error
	cmd := "basic"

	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}
	switch cmd {
	case "basic":
		err = BasicConfigurationHandling()
	default:
		err = fmt.Errorf("unknown example %q", cmd)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
