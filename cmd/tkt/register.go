package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/armando284/tkt/internal/db"
	"github.com/armando284/tkt/internal/logger"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register [path]",
	Short: "Register a project root for multi-project support",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var root string
		if len(args) > 0 {
			root = args[0]
		} else {
			var err error
			root, err = os.Getwd()
			if err != nil {
				return err
			}
		}

		absRoot, err := filepath.Abs(root)
		if err != nil {
			return err
		}

		name := filepath.Base(absRoot)

		_, err = db.DB.Exec(`
			INSERT OR REPLACE INTO projects (root_path, name)
			VALUES (?, ?)
		`, absRoot, name)

		if err != nil {
			return fmt.Errorf("failed to register project: %w", err)
		}

		logger.L.Info(fmt.Sprintf("✅ Project registered: %s", absRoot))
		logger.L.Info(fmt.Sprintf("   Name: %s", name))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
}
