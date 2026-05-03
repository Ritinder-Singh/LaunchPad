package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/mrmackaniel/launchpad/internal/config"
	tmpl "github.com/mrmackaniel/launchpad/internal/template"
	"github.com/spf13/cobra"
)

var previewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Serve the generated showcase page locally",
	RunE:  runPreview,
}

var previewPort int
var previewConfig string

func init() {
	previewCmd.Flags().IntVarP(&previewPort, "port", "p", 3000, "local port to serve on")
	previewCmd.Flags().StringVarP(&previewConfig, "config", "c", "showcase.config.json", "path to config file")
}

func runPreview(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(previewConfig)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	html, err := tmpl.Render(cfg)
	if err != nil {
		return fmt.Errorf("rendering template: %w", err)
	}

	addr := fmt.Sprintf(":%d", previewPort)
	fmt.Fprintf(os.Stdout, "Preview running at http://localhost%s\n", addr)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(html)
	})

	return http.ListenAndServe(addr, nil)
}
