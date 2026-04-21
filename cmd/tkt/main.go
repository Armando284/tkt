package main

import (
	"fmt"
	"os"

	"github.com/armando284/tkt/internal/config"
	"github.com/armando284/tkt/internal/db"
	"github.com/armando284/tkt/internal/logger"
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0-dev"
)

var rootCmd = &cobra.Command{
	Use:     "tkt",
	Short:   "tkt - Fast personal ticket & workflow tool",
	Long:    `tkt helps you turn TODO comments into tickets, manage git branches, and track time automatically.`,
	Version: version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Load(); err != nil {
			return err
		}
		if err := db.Init(); err != nil {
			logger.L.Error(fmt.Sprintf("failed to initialize database: %v", err))
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger.L.Info("🚀 tkt - Personal Ticket System")
		logger.L.Info(fmt.Sprintf("   Version %s", version))
		logger.L.Info("\nUse 'tkt --help' to see available commands.")
	},
}

func init() {
	logger.Init()
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(helloCmd)
	// New commands will be added here in future steps 
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		logger.L.Info(fmt.Sprintf("tkt version %s", version))
	},
}

var helloCmd = &cobra.Command{
	Use:   "hello [name]",
	Short: "Test command",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := "developer"
		if len(args) > 0 {
			name = args[0]
		}
		logger.L.Info(fmt.Sprintf("👋 Hello %s! tkt is ready.", name))
	},
}

func main() {
	defer db.Close()

	if err := rootCmd.Execute(); err != nil {
		logger.L.Error(fmt.Sprintf("Error: %v", err))
		os.Exit(1)
	}
}
