package main

import (
	"os"

	"github.com/truck8ai/battlefaeries-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
