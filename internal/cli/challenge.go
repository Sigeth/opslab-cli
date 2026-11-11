package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Sigeth/opslab-core-go/pkg/challenge"
)

func newChallengeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "challenge",
		Short: "Write, test and publish challenges",
	}
	cmd.AddCommand(newChallengeLintCommand())
	return cmd
}

func newChallengeLintCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "lint [DIR]",
		Short: "Validate a challenge package (default: current directory)",
		Long: "Validates " + challenge.ManifestFile + " and checks that every file it references exists.\n" +
			"It runs the same validation as the operator when the challenge is installed.",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) == 1 {
				dir = args[0]
			}
			return lint(cmd, dir)
		},
	}
}

var errInvalid = errors.New("invalid challenge package")

func lint(cmd *cobra.Command, dir string) error {
	files := os.DirFS(dir)
	pkg, err := challenge.Load(files)
	if err != nil {
		return fmt.Errorf("%s: %w", dir, err)
	}
	errs := challenge.Validate(pkg, files)
	out := cmd.OutOrStdout()
	if len(errs) == 0 {
		fmt.Fprintf(out, "%s %s: valid (%d machines, %d tasks)\n",
			pkg.Metadata.Name, pkg.Metadata.Version, len(pkg.Spec.Machines), len(pkg.Spec.Tasks))
		return nil
	}
	for _, e := range errs {
		fmt.Fprintf(out, "  %s\n", e.Error())
	}
	return fmt.Errorf("%w: %d problem(s) in %s", errInvalid, len(errs), dir)
}
