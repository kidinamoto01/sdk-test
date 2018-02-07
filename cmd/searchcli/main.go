package main

import (
	"os"
	"github.com/spf13/cobra"
	"github.com/tendermint/tmlibs/cli"

	"github.com/kidinamoto01/test/version"
	)

// entry point for this binary
var (
	ExplorerCli = &cobra.Command{
		Use:   "explorercli",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)


func main() {
	// disable sorting
	cobra.EnableCommandSorting = false



	ExplorerCli.AddCommand(
		//commands.InitCmd,
		//restServerCmd,
		//syncCmd,

		version.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(ExplorerCli, "EX", os.ExpandEnv("$HOME/.search-cli"))
	executor.Execute()
}