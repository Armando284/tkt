package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/armando284/tkt/internal/db"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit a ticket interactively",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		idStr := args[0]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fmt.Errorf("invalid ticket ID: %s", idStr)
		}

		// Obtener ticket actual
		var title, status, folder string
		err = db.DB.QueryRow(`
			SELECT title, status, folder 
			FROM tickets 
			WHERE id = ?
		`, id).Scan(&title, &status, &folder)

		if err != nil {
			return fmt.Errorf("ticket #%d not found", id)
		}

		fmt.Printf("Editing ticket #%d\n", id)
		fmt.Printf("Current title : %s\n", title)
		fmt.Printf("Current status: %s\n", status)
		fmt.Printf("Current folder: %s\n\n", folder)

		reader := bufio.NewReader(os.Stdin)

		// Editar título
		fmt.Print("New title (press Enter to keep current): ")
		newTitle, _ := reader.ReadString('\n')
		newTitle = strings.TrimSpace(newTitle)
		if newTitle == "" {
			newTitle = title
		}

		// Editar status
		fmt.Printf("New status (todo/in-progress/done) [current: %s]: ", status)
		newStatus, _ := reader.ReadString('\n')
		newStatus = strings.TrimSpace(newStatus)
		if newStatus == "" {
			newStatus = status
		}

		// Editar folder (opcional)
		fmt.Print("New folder (press Enter to keep current): ")
		newFolder, _ := reader.ReadString('\n')
		newFolder = strings.TrimSpace(newFolder)
		if newFolder == "" {
			newFolder = folder
		}

		_, err = db.DB.Exec(`
			UPDATE tickets 
			SET title = ?, status = ?, folder = ?
			WHERE id = ?
		`, newTitle, newStatus, newFolder, id)

		if err != nil {
			return fmt.Errorf("failed to update ticket: %w", err)
		}

		fmt.Printf("\n✅ Ticket #%d updated successfully!\n", id)
		fmt.Printf("   Title : %s\n", newTitle)
		fmt.Printf("   Status: %s\n", newStatus)
		fmt.Printf("   Folder: %s\n", newFolder)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}