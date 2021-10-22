package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var flagDataDir = "datadir"

func main() {
	var tbbCmd = &cobra.Command{
		Use:   "tbb",
		Short: "The Blockchain Bar CLI",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	tbbCmd.AddCommand(versionCmd)
	tbbCmd.AddCommand(balancesCmd())
	tbbCmd.AddCommand(txCmd())

	err := tbbCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func addDefaultRequiredFlags(cmd *cobra.Command) {
	cmd.Flags().String(
		flagDataDir,
		"",
		"Absolute path where all data will/is stored",
	)
	cmd.MarkFlagRequired(flagDataDir)
}
