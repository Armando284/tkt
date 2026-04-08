package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// version will be injected at build time via ldflags
	version = "0.1.0-dev"
)

var rootCmd = &cobra.Command{
	Use:   "tkt",
	Short: "tkt - Fast personal ticket & workflow tool",
	Long: `tkt is a lightweight and fast Go CLI tool designed for developers.

It turns TODO/FIXME comments into tickets, manages git branches per ticket,
tracks time automatically when using Neovim, and supports multiple projects.`,
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 tkt - Personal Ticket System")
		fmt.Printf("   Version %s\n", version)
		fmt.Println("\nRun 'tkt --help' for available commands.")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(helloCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of tkt",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("tkt version %s\n", version)
	},
}

var helloCmd = &cobra.Command{
	Use:   "hello [name]",
	Short: "Say hello (test command)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := "developer"
		if len(args) > 0 {
			name = args[0]
		}
		fmt.Printf("👋 Hello %s!\n", name)
		fmt.Println("Welcome to tkt — your terminal ticket workflow tool.")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
