package cmd

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mrmackaniel/launchpad/internal/config"
	tmpl "github.com/mrmackaniel/launchpad/internal/template"
	"github.com/spf13/cobra"
)

var previewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Render and serve the showcase page locally",
	RunE:  runPreview,
}

var previewPort int
var previewConfig string

func init() {
	previewCmd.Flags().IntVarP(&previewPort, "port", "p", 3000, "local port to serve on")
	previewCmd.Flags().StringVarP(&previewConfig, "config", "c", config.DefaultConfigPath, "path to config file")
}

func runPreview(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(previewConfig)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Prefer dist screenshots (post-deploy), fall back to source screenshots
	screenshotSrc := filepath.Join(".showcase", "dist", "screenshots")
	if entries, _ := os.ReadDir(screenshotSrc); len(entries) == 0 {
		screenshotSrc = filepath.Join(".showcase", "screenshots")
	}

	var screenshotRels []string
	if entries, err := os.ReadDir(screenshotSrc); err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				screenshotRels = append(screenshotRels, "screenshots/"+e.Name())
			}
		}
	}

	html, err := tmpl.Render(cfg, screenshotRels)
	if err != nil {
		return fmt.Errorf("rendering template: %w", err)
	}

	addr := fmt.Sprintf(":%d", previewPort)
	fmt.Fprintf(os.Stdout, "Preview at http://localhost%s\n", addr)

	mux := http.NewServeMux()
	mux.Handle("/screenshots/", http.StripPrefix("/screenshots/", http.FileServer(http.Dir(screenshotSrc))))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(html)
	})

	return http.ListenAndServe(addr, mux)
}
