package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version will be set at build time
	version = "0.1.0-dev"
)

var rootCmd = &cobra.Command{
	Use:     "tkt",
	Short:   "tkt - Fast personal ticket & workflow tool for developers",
	Long: `tkt is a lightweight Go CLI that helps developers manage tickets
directly from their codebase.

Features:
- Scan TODO/FIXME comments → create tickets
- Git branch management per ticket
- Automatic time tracking with Neovim
- Multi-project support from your home directory
- Efficient real-time file watcher`,
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 tkt - Personal Ticket System")
		fmt.Println("   Fast • Lightweight • Terminal-first")
		fmt.Println("\nRun 'tkt --help' to see all commands.")
	},
}

func init() {
	// Agregamos el comando version de forma explícita
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(helloCmd)
}

// versionCmd muestra la versión
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of tkt",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("tkt version %s\n", version)
	},
}

// helloCmd es nuestro primer comando de prueba
var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Say hello to tkt",
	Run: func(cmd *cobra.Command, args []string) {
		name := "developer"
		if len(args) > 0 {
			name = args[0]
		}
		fmt.Printf("👋 Hello %s! Welcome to tkt.\n", name)
		fmt.Println("Your personal ticket workflow tool is ready.")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
