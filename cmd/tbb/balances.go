package main

import (
	"fmt"
	"os"
	"the-blockchain-bar/database"

	"github.com/spf13/cobra"
)

func balancesCmd() *cobra.Command {
	var balancesCmd = &cobra.Command{
		Use:   "balances",
		Short: "Interact with balances (list...).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {},
	}

	balancesCmd.AddCommand(balancesListCmd())

	return balancesCmd
}

func incorrectUsageErr() error {
	return fmt.Errorf("incorrect usage")
}

func balancesListCmd() *cobra.Command {
	var balancesListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all balances",
		Run: func(cmd *cobra.Command, args []string) {
			dataDir, _ := cmd.Flags().GetString(flagDataDir)
			state, err := database.NewStateFromDisk(dataDir)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer state.Close()

			fmt.Printf("Account balances at: %x\n", state.LatestSnapshot())
			fmt.Println("---------------------")
			for account, balance := range state.Balances {
				fmt.Printf(fmt.Sprintf("%s, %d", account, balance) + "\n")
			}
		},
	}

	addDefaultRequiredFlags(balancesListCmd)

	return balancesListCmd
}
