package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/armando284/tkt/internal/db"
	"github.com/armando284/tkt/internal/logger"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [id]",
	Short: "Start working on a ticket",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.L.Debug("start command started")

		var ticketID int
		var err error

		// ==================== MODO CON ID O INTERACTIVO ====================
		if len(args) > 0 {
			logger.L.Debug(fmt.Sprintf("Modo directo con ID = %s", args[0]))
			ticketID, err = strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("ID inválido: %s", args[0])
			}
		} else {
			logger.L.Debug("Modo interactivo - mostrando lista")
			rows, err := db.DB.Query("SELECT id, status, title FROM tickets ORDER BY id DESC")
			if err != nil {
				return fmt.Errorf("error consultando tickets: %w", err)
			}
			defer rows.Close()

			logger.L.Info("ID   | Status       | Title")
			logger.L.Info("-----|--------------|-----------------------------")

			for rows.Next() {
				var id int
				var status, title string
				if err := rows.Scan(&id, &status, &title); err == nil {
					logger.L.Info(fmt.Sprintf("%-4d | %-12s | %s", id, status, title))
				}
			}

			logger.L.Info("\nElige el ID del ticket: ")
			_, err = fmt.Scanln(&ticketID)
			if err != nil {
				return fmt.Errorf("entrada inválida")
			}
			logger.L.Debug(fmt.Sprintf("Usuario seleccionó ID = %d", ticketID))
		}

		// ==================== CONSULTA SIMPLE AL TICKET ====================
		logger.L.Debug(fmt.Sprintf("Buscando ticket ID = %d en la BD...", ticketID))

		var title, projectRoot string
		var folder, branch sql.NullString

		err = db.DB.QueryRow(`
			SELECT title, folder, branch, project_root 
			FROM tickets 
			WHERE id = ?
		`, ticketID).Scan(&title, &folder, &branch, &projectRoot)

		if err != nil {
			logger.L.Debug(fmt.Sprintf("QueryRow falló con error: %v", err))
			return fmt.Errorf("❌ Ticket #%d NO encontrado", ticketID)
		}

		logger.L.Debug("Ticket encontrado correctamente")
		logger.L.Debug(fmt.Sprintf("title       = '%s'", title))
		logger.L.Debug(fmt.Sprintf("folder      = '%s'", folder.String))
		logger.L.Debug(fmt.Sprintf("branch      = '%s'", branch.String))
		logger.L.Debug(fmt.Sprintf("project_root = '%s'", projectRoot))

		// Determinar carpeta final
		finalFolder := folder.String
		if finalFolder == "" {
			finalFolder = projectRoot
			logger.L.Debug("Usando project_root como folder")
		}
		if finalFolder == "" {
			finalFolder = "."
			logger.L.Debug("Usando '.' como folder")
		}

		// Generar rama si no tiene
		if branch.String == "" {
			branchStr := fmt.Sprintf("feature/tkt-%04d-%s", ticketID, kebabCase(title))
			logger.L.Debug(fmt.Sprintf("Generando rama: %s", branchStr))
			_, _ = db.DB.Exec("UPDATE tickets SET branch = ? WHERE id = ?", branchStr, ticketID)
			branch.String = branchStr // actualizar localmente
		}

		// Actualizar estado a in-progress
		logger.L.Debug("Actualizando ticket a 'in-progress'...")
		_, err = db.DB.Exec(`
			UPDATE tickets 
			SET status = 'in-progress', 
			    current_start_ts = datetime('now', 'localtime')
			WHERE id = ?
		`, ticketID)
		if err != nil {
			logger.L.Debug(fmt.Sprintf("Error al actualizar estado: %v", err))
		} else {
			logger.L.Debug("Estado actualizado correctamente")
		}

		// Salida para el wrapper bash
		logger.L.Info(fmt.Sprintf("CD:%s", finalFolder))
		logger.L.Info(fmt.Sprintf("TICKET_ID:%s", strconv.Itoa(ticketID)))
		logger.L.Info(fmt.Sprintf("BRANCH:%s", branch.String))

		logger.L.Info(fmt.Sprintf("✅ Iniciando ticket #%d → %s", ticketID, title))
		return nil
	},
}

func kebabCase(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, ":", "")
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", "")
	return s
}

func init() {
	rootCmd.AddCommand(startCmd)
}