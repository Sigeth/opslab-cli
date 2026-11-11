// Command opslab installs and operates the opslab platform, and helps authors
// write, test and publish challenges.
package main

import (
	"os"

	"github.com/Sigeth/opslab-cli/internal/cli"
)

func main() {
	if err := cli.NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
