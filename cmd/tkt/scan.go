package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/armando284/tkt/internal/db"
	"github.com/armando284/tkt/internal/utils"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan all registered projects for TODO/FIXME comments and create tickets",
	RunE: func(cmd *cobra.Command, args []string) error {
		rows, err := db.DB.Query("SELECT root_path, name FROM projects")
		if err != nil {
			return fmt.Errorf("failed to get projects: %w", err)
		}
		defer rows.Close()

		var projects []struct {
			root string
			name string
		}

		for rows.Next() {
			var root, name string
			if err := rows.Scan(&root, &name); err != nil {
				continue
			}
			projects = append(projects, struct {
				root string
				name string
			}{root, name})
		}

		if len(projects) == 0 {
			fmt.Println("❌ No projects registered yet. Use 'tkt register' first.")
			return nil
		}

		fmt.Printf("🔍 Scanning %d project(s) for TODO comments...\n\n", len(projects))

		todoRegex := regexp.MustCompile(`(?im)(?:^|[\s])(?:\/\/|/\*|#)\s*(TODO|FIXME|HACK):\s*(.+?)(?:\s*\*/|$)`)
		newTickets := 0

		for _, proj := range projects {
			fmt.Printf("📁 Scanning: %s\n", proj.name)

			err := filepath.WalkDir(proj.root, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return nil
				}

				if d.IsDir() {
					name := d.Name()
					if name == ".git" || name == "node_modules" || name == ".venv" || name == "venv" ||
						name == "dist" || name == "build" || name == "target" {
						return filepath.SkipDir
					}
					return nil
				}

				ext := strings.ToLower(filepath.Ext(path))
				if !isCodeFile(ext) {
					return nil
				}

				utils.Dev.Debugf("scanning file: %s", path)

				content, err := os.ReadFile(path)
				if err != nil {
					return nil
				}

				matches := todoRegex.FindAllStringSubmatch(string(content), -1)
				for _, match := range matches {
					if len(match) < 3 {
						continue
					}
					title := strings.TrimSpace(match[2])
					if title == "" || len(title) > 200 {
						continue
					}

					utils.Dev.Debugf("found TODO: %s in %s", title, path)

					_, err := db.DB.Exec(`
						INSERT INTO tickets (title, folder, project_root, status)
						VALUES (?, ?, ?, 'todo')
						ON CONFLICT(title, project_root) DO NOTHING
					`, title, filepath.Dir(path), proj.root)

					if err != nil {
						utils.Dev.Errorf("failed to insert ticket %q: %v", title, err)
						continue
					}

					utils.Dev.Debugf("inserted ticket: %s", title)
					newTickets++
					fmt.Printf("  ✅ Created ticket: %s\n", title)
				}
				return nil
			})

			if err != nil {
				fmt.Printf("  ⚠️  Error scanning %s: %v\n", proj.name, err)
			}
		}

		fmt.Printf("\n🎉 Scan completed! %d new tickets created.\n", newTickets)
		return nil
	},
}

func isCodeFile(ext string) bool {
	codeExts := map[string]bool{
		".go": true, ".js": true, ".ts": true, ".tsx": true, ".jsx": true,
		".py": true, ".rs": true, ".java": true, ".c": true, ".cpp": true,
		".cs": true, ".php": true, ".rb": true, ".md": true, ".txt": true,
	}
	return codeExts[ext]
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
