package main

import (
	"fmt"
	"github.com/armando284/tkt/internal/db"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tickets across all registered projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		rows, err := db.DB.Query(`
			SELECT t.id, t.title, t.status, t.folder, p.name as project_name
			FROM tickets t
			JOIN projects p ON t.project_root = p.root_path
			ORDER BY t.id DESC
		`)
		if err != nil {
			return fmt.Errorf("query failed: %w", err)
		}
		defer rows.Close()

		fmt.Println("ID   | Status       | Project          | Title")
		fmt.Println("-----|--------------|------------------|-----------------------------------")

		count := 0
		for rows.Next() {
			var id int
			var title, status, folder, project string
			if err := rows.Scan(&id, &title, &status, &folder, &project); err != nil {
				continue
			}
			fmt.Printf("%-4d | %-12s | %-16s | %s\n", id, status, project, title)
			count++
		}

		if count == 0 {
			fmt.Println("No tickets found yet. Run 'tkt scan' after registering projects.")
		} else {
			fmt.Printf("\nTotal tickets: %d\n", count)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
