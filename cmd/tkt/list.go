package main

import (
	"fmt"

	"github.com/armando284/tkt/internal/db"
	"github.com/armando284/tkt/internal/logger"
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

		logger.L.Info("ID   | Status       | Project          | Title")
		logger.L.Info("-----|--------------|------------------|-----------------------------------")

		count := 0
		for rows.Next() {
			var id int
			var title, status, folder, project string
			if err := rows.Scan(&id, &title, &status, &folder, &project); err != nil {
				continue
			}
			logger.L.Info(fmt.Sprintf("%-4d | %-12s | %-16s | %s", id, status, project, title))
			count++
		}

		if count == 0 {
			logger.L.Info("No tickets found yet. Run 'tkt scan' after registering projects.")
		} else {
			logger.L.Info(fmt.Sprintf("Total tickets: %d", count))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
