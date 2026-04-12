package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/armando284/tkt/internal/db"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [title]",
	Short: "Create a new ticket manually",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var title string

		if len(args) > 0 {
			title = args[0]
		} else {
			// Modo interactivo
			fmt.Print("Enter ticket title: ")
			reader := bufio.NewReader(os.Stdin)
			title, _ = reader.ReadString('\n')
			title = strings.TrimSpace(title)
		}

		if title == "" {
			return fmt.Errorf("title cannot be empty")
		}

		// Usamos el directorio actual como folder por defecto
		folder, _ := os.Getwd()

		_, err := db.DB.Exec(`
			INSERT INTO tickets (title, folder, project_root, status)
			VALUES (?, ?, ?, 'todo')
		`, title, folder, folder) // por ahora usamos folder como project_root

		if err != nil {
			return fmt.Errorf("failed to create ticket: %w", err)
		}

		fmt.Printf("✅ Ticket created successfully!\n")
		fmt.Printf("   Title : %s\n", title)
		fmt.Printf("   Status: todo\n")
		fmt.Printf("   Folder: %s\n", folder)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}