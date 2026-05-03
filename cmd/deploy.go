package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mrmackaniel/launchpad/internal/config"
	"github.com/mrmackaniel/launchpad/internal/deploy"
	tmpl "github.com/mrmackaniel/launchpad/internal/template"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Build and deploy the showcase page to Cloudflare Pages",
	RunE:  runDeploy,
}

var configPath string

func init() {
	deployCmd.Flags().StringVarP(&configPath, "config", "c", "showcase.config.json", "path to config file")
}

func runDeploy(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	fmt.Println("Rendering template...")
	html, err := tmpl.Render(cfg)
	if err != nil {
		return fmt.Errorf("rendering template: %w", err)
	}

	distDir, err := os.MkdirTemp("", "launchpad-dist-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(distDir)

	if err := os.WriteFile(filepath.Join(distDir, "index.html"), html, 0644); err != nil {
		return err
	}

	fmt.Printf("Deploying %s to %s.%s...\n", cfg.Name, cfg.Deploy.Subdomain, cfg.Deploy.Domain)

	deployer, err := deploy.New(cfg.Deploy.Provider)
	if err != nil {
		return err
	}

	url, err := deployer.Deploy(cfg, distDir)
	if err != nil {
		return fmt.Errorf("deploy failed: %w", err)
	}

	fmt.Printf("\nLive at: %s\n", url)
	return nil
}
