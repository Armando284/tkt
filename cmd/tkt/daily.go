package main

import (
	"fmt"
	"time"

	"github.com/armando284/tkt/internal/db"
	"github.com/armando284/tkt/internal/logger"
	"github.com/spf13/cobra"
)

var dailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "Show time spent on tickets for a specific day",
	Long: `Muestra el tiempo dedicado a tickets en un día específico.
Sin argumentos muestra el día de hoy.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dateStr, _ := cmd.Flags().GetString("date")

		var targetDate string
		if dateStr != "" {
			// Validar formato de fecha
			_, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return fmt.Errorf("formato de fecha inválido. Usa YYYY-MM-DD (ej: 2026-04-18)")
			}
			targetDate = dateStr
		} else {
			// Hoy por defecto
			targetDate = time.Now().Format("2006-01-02")
		}

		logger.L.Info(fmt.Sprintf("=== Reporte de tiempo - %s ===", targetDate))

		rows, err := db.DB.Query(`
			SELECT 
				t.id,
				t.title,
				t.status,
				SUM(s.duration) as total_seconds,
				COUNT(s.id) as session_count,
				MIN(s.start_ts) as first_start,
				MAX(s.end_ts) as last_end
			FROM sessions s
			JOIN tickets t ON s.ticket_id = t.id
			WHERE s.start_ts LIKE ? || '%%'
			GROUP BY t.id, t.title, t.status
			ORDER BY total_seconds DESC
		`, targetDate)

		if err != nil {
			return fmt.Errorf("error consultando sesiones: %w", err)
		}
		defer rows.Close()

		var totalDaySeconds int
		hasResults := false

		logger.L.Info("ID   | Status       | Tiempo     | Sesiones | Título")
		logger.L.Info("-----|--------------|------------|----------|-----------------------------------")

		for rows.Next() {
			hasResults = true
			var id, totalSec, sessions int
			var title, status, firstStart, lastEnd string

			if err := rows.Scan(&id, &title, &status, &totalSec, &sessions, &firstStart, &lastEnd); err != nil {
				continue
			}

			totalDaySeconds += totalSec

			hours := totalSec / 3600
			minutes := (totalSec % 3600) / 60

			timeStr := fmt.Sprintf("%dh %dm", hours, minutes)
			if hours == 0 {
				timeStr = fmt.Sprintf("%dm", minutes)
			}

			logger.L.Info(fmt.Sprintf("%-4d | %-12s | %-10s | %-8d | %s",
				id, status, timeStr, sessions, title))
		}

		if !hasResults {
			logger.L.Info(fmt.Sprintf("No se registró tiempo el día %s.", targetDate))
			return nil
		}

		// Total del día
		totalHours := totalDaySeconds / 3600
		totalMinutes := (totalDaySeconds % 3600) / 60

		logger.L.Info("-------------------------------------------------------------")
		logger.L.Info(fmt.Sprintf("Total del día: %dh %dm (%d minutos)", 
			totalHours, totalMinutes, totalDaySeconds/60))

		return nil
	},
}

func init() {
	dailyCmd.Flags().String("date", "", "Fecha en formato YYYY-MM-DD (ej: 2026-04-18)")
	rootCmd.AddCommand(dailyCmd)
}