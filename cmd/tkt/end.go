package main

import (
	"fmt"
	"strconv"

	"github.com/armando284/tkt/internal/db"
	"github.com/armando284/tkt/internal/logger"
	"github.com/spf13/cobra"
)

var endCmd = &cobra.Command{
	Use:   "end --id <id>",
	Short: "End current work session and save time",
	RunE: func(cmd *cobra.Command, args []string) error {
		idStr, _ := cmd.Flags().GetString("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fmt.Errorf("invalid ticket ID")
		}

		var startTs string
		err = db.DB.QueryRow("SELECT current_start_ts FROM tickets WHERE id = ?", id).Scan(&startTs)
		if err != nil || startTs == "" {
			logger.L.Info("No active session found for this ticket.")
			return nil
		}

		_, err = db.DB.Exec(`
			INSERT INTO sessions (ticket_id, start_ts, end_ts, duration)
			SELECT id, current_start_ts, datetime('now', 'localtime'),
			       strftime('%s', 'now') - strftime('%s', current_start_ts)
			FROM tickets WHERE id = ?
		`, id)

		if err != nil {
			return err
		}

		_, err = db.DB.Exec("UPDATE tickets SET current_start_ts = NULL WHERE id = ?", id)
		if err != nil {
			return err
		}

		logger.L.Info(fmt.Sprintf("✅ Session ended for ticket #%d. Time saved.", id))
		return nil
	},
}

func init() {
	endCmd.Flags().String("id", "", "Ticket ID")
	endCmd.MarkFlagRequired("id")
	rootCmd.AddCommand(endCmd)
}