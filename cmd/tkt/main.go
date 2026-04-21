package main

import (
	"fmt"
	"os"

	"github.com/armando284/tkt/internal/config"
	"github.com/armando284/tkt/internal/db"
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
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 tkt - Personal Ticket System")
		fmt.Printf("   Version %s\n", version)
		fmt.Println("\nUse 'tkt --help' to see available commands.")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(helloCmd)
	// New commands will be added here in future steps 
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("tkt version %s\n", version)
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
		fmt.Printf("👋 Hello %s! tkt is ready.\n", name)
	},
}

func main() {
	defer db.Close()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
