package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/armando284/tkt/internal/db"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <id>[,<id2>,...]",
	Short: "Delete one or multiple tickets by ID",
	Example: `  tkt delete 42
  tkt delete 5,7,12
  tkt delete 3,8`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		idStr := strings.TrimSpace(args[0])
		if idStr == "" {
			return fmt.Errorf("no ticket IDs provided")
		}

		// Parsear IDs separados por coma
		idParts := strings.Split(idStr, ",")
		var ids []int

		for _, part := range idParts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			id, err := strconv.Atoi(part)
			if err != nil {
				return fmt.Errorf("invalid ticket ID: %s", part)
			}
			if id <= 0 {
				return fmt.Errorf("invalid ticket ID: %d (must be positive)", id)
			}
			ids = append(ids, id)
		}

		if len(ids) == 0 {
			return fmt.Errorf("no valid ticket IDs provided")
		}

		// Obtener información de los tickets antes de borrar
		type ticketInfo struct {
			ID    int
			Title string
			Status string
		}

		var tickets []ticketInfo

		for _, id := range ids {
			var title, status string
			err := db.DB.QueryRow(`
				SELECT title, status 
				FROM tickets 
				WHERE id = ?
			`, id).Scan(&title, &status)

			if err != nil {
				fmt.Printf("⚠️  Ticket #%d not found, skipping...\n", id)
				continue
			}

			tickets = append(tickets, ticketInfo{ID: id, Title: title, Status: status})
		}

		if len(tickets) == 0 {
			fmt.Println("❌ No tickets found to delete.")
			return nil
		}

		// Mostrar tickets que se van a eliminar
		fmt.Println("The following tickets will be deleted:")
		fmt.Println("--------------------------------------------------")
		for _, t := range tickets {
			fmt.Printf("#%-4d | %-12s | %s\n", t.ID, t.Status, t.Title)
		}
		fmt.Println("--------------------------------------------------")

		// Confirmación general
		if len(tickets) > 1 {
			fmt.Printf("Are you sure you want to delete these %d tickets? (y/N): ", len(tickets))
		} else {
			fmt.Print("Are you sure you want to delete this ticket? (y/N): ")
		}

		var confirm string
		fmt.Scanln(&confirm)

		if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
			fmt.Println("Delete operation cancelled.")
			return nil
		}

		// Borrar uno por uno con confirmación individual si son varios
		deleted := 0
		for _, t := range tickets {
			if len(tickets) > 1 {
				fmt.Printf("\nDelete ticket #%d: %s ? (y/N): ", t.ID, t.Title)
				var singleConfirm string
				fmt.Scanln(&singleConfirm)

				if strings.ToLower(strings.TrimSpace(singleConfirm)) != "y" {
					fmt.Printf("Skipping ticket #%d\n", t.ID)
					continue
				}
			}

			result, err := db.DB.Exec("DELETE FROM tickets WHERE id = ?", t.ID)
			if err != nil {
				fmt.Printf("❌ Failed to delete ticket #%d: %v\n", t.ID, err)
				continue
			}

			if rows, _ := result.RowsAffected(); rows > 0 {
				fmt.Printf("✅ Deleted ticket #%d\n", t.ID)
				deleted++
			}
		}

		fmt.Printf("\nDeletion completed. %d ticket(s) deleted.\n", deleted)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}