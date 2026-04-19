package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/armando284/tkt/internal/db"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [id]",
	Short: "Start working on a ticket",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("=== DEBUG: start command started ===")

		var ticketID int
		var err error

		// ==================== MODO CON ID O INTERACTIVO ====================
		if len(args) > 0 {
			fmt.Printf("DEBUG: Modo directo con ID = %s\n", args[0])
			ticketID, err = strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("ID inválido: %s", args[0])
			}
		} else {
			fmt.Println("DEBUG: Modo interactivo - mostrando lista")
			rows, err := db.DB.Query("SELECT id, status, title FROM tickets ORDER BY id DESC")
			if err != nil {
				return fmt.Errorf("error consultando tickets: %w", err)
			}
			defer rows.Close()

			fmt.Println("ID   | Status       | Title")
			fmt.Println("-----|--------------|-----------------------------")

			for rows.Next() {
				var id int
				var status, title string
				if err := rows.Scan(&id, &status, &title); err == nil {
					fmt.Printf("%-4d | %-12s | %s\n", id, status, title)
				}
			}

			fmt.Print("\nElige el ID del ticket: ")
			_, err = fmt.Scanln(&ticketID)
			if err != nil {
				return fmt.Errorf("entrada inválida")
			}
			fmt.Printf("DEBUG: Usuario seleccionó ID = %d\n", ticketID)
		}

		// ==================== CONSULTA SIMPLE AL TICKET ====================
		fmt.Printf("DEBUG: Buscando ticket ID = %d en la BD...\n", ticketID)

		var title, folder, branch, projectRoot string

		err = db.DB.QueryRow(`
			SELECT title, folder, branch, project_root 
			FROM tickets 
			WHERE id = ?
		`, ticketID).Scan(&title, &folder, &branch, &projectRoot)

		if err != nil {
			fmt.Printf("DEBUG: QueryRow falló con error: %v\n", err)
			return fmt.Errorf("❌ Ticket #%d NO encontrado", ticketID)
		}

		fmt.Println("DEBUG: Ticket encontrado correctamente")
		fmt.Printf("DEBUG: title       = '%s'\n", title)
		fmt.Printf("DEBUG: folder      = '%s'\n", folder)
		fmt.Printf("DEBUG: branch      = '%s'\n", branch)
		fmt.Printf("DEBUG: project_root = '%s'\n", projectRoot)

		// Determinar carpeta final
		finalFolder := folder
		if finalFolder == "" {
			finalFolder = projectRoot
			fmt.Println("DEBUG: Usando project_root como folder")
		}
		if finalFolder == "" {
			finalFolder = "."
			fmt.Println("DEBUG: Usando '.' como folder")
		}

		// Generar rama si no tiene
		if branch == "" {
			branch = fmt.Sprintf("feature/tkt-%04d-%s", ticketID, kebabCase(title))
			fmt.Printf("DEBUG: Generando rama: %s\n", branch)
			_, _ = db.DB.Exec("UPDATE tickets SET branch = ? WHERE id = ?", branch, ticketID)
		}

		// Actualizar estado a in-progress
		fmt.Println("DEBUG: Actualizando ticket a 'in-progress'...")
		_, err = db.DB.Exec(`
			UPDATE tickets 
			SET status = 'in-progress', 
			    current_start_ts = datetime('now', 'localtime')
			WHERE id = ?
		`, ticketID)
		if err != nil {
			fmt.Printf("DEBUG: Error al actualizar estado: %v\n", err)
		} else {
			fmt.Println("DEBUG: Estado actualizado correctamente")
		}

		// Salida para el wrapper bash
		fmt.Println("CD:" + finalFolder)
		fmt.Println("TICKET_ID:" + strconv.Itoa(ticketID))
		fmt.Println("BRANCH:" + branch)

		fmt.Printf("✅ Iniciando ticket #%d → %s\n", ticketID, title)
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