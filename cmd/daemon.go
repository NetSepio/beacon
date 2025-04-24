package cmd

import (
	"github.com/NetSepio/beacon/core"
	"github.com/spf13/cobra"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run the Erebrus node as a daemon",
	Run: func(cmd *cobra.Command, args []string) {
		core.RunBeaconNode()
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)
} 