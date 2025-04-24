package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/NetSepio/beacon/core"
	"github.com/NetSepio/beacon/util"
	"github.com/spf13/cobra"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/patrickmn/go-cache"
	"github.com/NetSepio/beacon/p2p"
	helmet "github.com/danielkov/gin-helmet"
	log "github.com/sirupsen/logrus"
)

var (
	// ANSI color codes
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

var rootCmd = &cobra.Command{
	Use:   "erebrus",
	Short: "Erebrus is a decentralized VPN node",
	Long: `Erebrus is a decentralized VPN node that provides secure and private internet access.
Complete documentation is available at https://erebrus.io`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Erebrus",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\n%s%s%s\n", colorYellow, "====================================", colorReset)
		fmt.Printf("%s📦 Erebrus Version%s\n", colorGreen, colorReset)
		fmt.Printf("%s%s%s\n", colorYellow, "====================================", colorReset)
		fmt.Printf("%s🔖 Version:%s %s\n", colorCyan, colorReset, util.Version)
		fmt.Printf("%s🔖 Code Hash:%s %s\n", colorCyan, colorReset, util.CodeHash)
		fmt.Printf("%s🔖 Go Version:%s %s\n", colorCyan, colorReset, util.GoVersion)
		fmt.Printf("%s%s%s\n\n", colorYellow, "====================================", colorReset)
	},
}

var deactivateCmd = &cobra.Command{
	Use:   "deactivate",
	Short: "Deactivate the Erebrus node",
	Run: func(cmd *cobra.Command, args []string) {
		if err := core.DeactivateNode(); err != nil {
			fmt.Printf("\n%s❌ Error: %s%s\n", colorRed, err.Error(), colorReset)
			os.Exit(1)
		}
		fmt.Printf("%s✅ Node successfully deactivated%s\n", colorGreen, colorReset)
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the current status of the Erebrus node",
	Run: func(cmd *cobra.Command, args []string) {
		status, err := core.GetNodeStatus()
		if err != nil {
			fmt.Printf("\n%s❌ Error: %s%s\n", colorRed, err.Error(), colorReset)
			os.Exit(1)
		}
		fmt.Printf("\n%s%s%s\n", colorYellow, "====================================", colorReset)
		fmt.Printf("%s📊 Node Status%s\n", colorGreen, colorReset)
		fmt.Printf("%s%s%s\n", colorYellow, "====================================", colorReset)
		fmt.Printf("%s🆔 Node ID:%s %s\n", colorCyan, colorReset, status.ID)
		fmt.Printf("%s📛 Name:%s %s\n", colorCyan, colorReset, status.Name)
		fmt.Printf("%s📝 Spec:%s %s\n", colorCyan, colorReset, status.Spec)
		fmt.Printf("%s⚙️  Config:%s %s\n", colorCyan, colorReset, status.Config)
		fmt.Printf("%s🌐 IP Address:%s %s\n", colorCyan, colorReset, status.IPAddress)
		fmt.Printf("%s🗺  Region:%s %s\n", colorCyan, colorReset, status.Region)
		fmt.Printf("%s%s%s\n\n", colorYellow, "====================================", colorReset)
	},
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run the Erebrus node as a daemon",
	Run: func(cmd *cobra.Command, args []string) {
		core.RunBeaconNode()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(deactivateCmd)
	rootCmd.AddCommand(daemonCmd)
}

func Execute() error {
	return rootCmd.Execute()
} 