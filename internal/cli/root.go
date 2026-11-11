// Package cli implements the opslab commands.
package cli

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"

	"github.com/Sigeth/opslab-core-go/pkg/project"
)

// version is set at build time: -ldflags "-X github.com/Sigeth/opslab-cli/internal/cli.version=v0.1.0".
var version = ""

// NewRootCommand builds the command tree.
func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:           project.Name,
		Short:         "Install opslab and write challenges for it",
		SilenceUsage:  true,
		SilenceErrors: false,
	}
	root.AddCommand(newVersionCommand(), newChallengeCommand())
	return root
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the CLI version",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Fprintln(cmd.OutOrStdout(), currentVersion())
		},
	}
}

func currentVersion() string {
	if version != "" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "dev"
}
