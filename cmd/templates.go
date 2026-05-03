package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Manage showcase templates",
}

var templatesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available built-in templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		templates := []struct {
			name        string
			description string
		}{
			{"minimal", "Clean, whitespace-heavy, typography-focused"},
			{"bold", "Dark background, large hero, accent color"},
			{"developer", "Terminal/code aesthetic, monospace font"},
		}
		for _, t := range templates {
			fmt.Printf("  %-12s %s\n", t.name, t.description)
		}
		return nil
	},
}

func init() {
	templatesCmd.AddCommand(templatesListCmd)
}
