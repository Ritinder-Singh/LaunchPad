package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mrmackaniel/launchpad/internal/capture"
	"github.com/mrmackaniel/launchpad/internal/config"
	"github.com/mrmackaniel/launchpad/internal/deploy"
	"github.com/mrmackaniel/launchpad/internal/gitops"
	"github.com/mrmackaniel/launchpad/internal/project"
	tmpl "github.com/mrmackaniel/launchpad/internal/template"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Build, screenshot, and deploy the showcase page",
	RunE:  runDeploy,
}

var configPath string
var skipGitPush bool
var skipBuild bool

func init() {
	deployCmd.Flags().StringVarP(&configPath, "config", "c", config.DefaultConfigPath, "path to config file")
	deployCmd.Flags().BoolVar(&skipGitPush, "skip-git-push", false, "skip committing showcase to the showcase branch")
	deployCmd.Flags().BoolVar(&skipBuild, "skip-build", false, "skip building the project (use existing build output)")
}

func runDeploy(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	distDir := filepath.Join(".showcase", "dist")
	if err := os.MkdirAll(distDir, 0755); err != nil {
		return err
	}

	// 1. Detect and build project
	proj := project.Detect(".")
	proj.Apply(cfg.Build.Command, cfg.Build.StartCommand, cfg.Build.Port)

	if !skipBuild {
		if err := proj.Build(); err != nil {
			return fmt.Errorf("building project: %w", err)
		}
	}

	// 2. Capture screenshots
	screenshotDir := cfg.Screenshots.Dir
	absScreenshotDir, _ := filepath.Abs(screenshotDir)

	screenshots, err := capture.Screenshots(
		screenshotDir,
		cfg.Screenshots.AutoCapture,
		proj.Start,
		proj.WaitReady,
		proj.Stop,
	)
	if err != nil {
		return fmt.Errorf("capturing screenshots: %w", err)
	}

	// Copy screenshots into dist and build relative paths for templates
	var screenshotRels []string
	if len(screenshots) > 0 {
		screenshotDistDir := filepath.Join(distDir, "screenshots")
		if err := os.MkdirAll(screenshotDistDir, 0755); err != nil {
			return err
		}
		for _, src := range screenshots {
			base := filepath.Base(src)
			dst := filepath.Join(screenshotDistDir, base)
			if err := copyFile(src, dst); err != nil {
				return fmt.Errorf("copying screenshot: %w", err)
			}
			screenshotRels = append(screenshotRels, "screenshots/"+base)
		}
		_ = absScreenshotDir
	}

	// 3. Render template
	fmt.Println("Rendering template...")
	html, err := tmpl.Render(cfg, screenshotRels)
	if err != nil {
		return fmt.Errorf("rendering template: %w", err)
	}

	if err := os.WriteFile(filepath.Join(distDir, "index.html"), html, 0644); err != nil {
		return err
	}

	// 4. Git: push showcase branch
	if !skipGitPush {
		if err := gitops.PushShowcase(distDir); err != nil {
			fmt.Printf("Warning: git push failed: %v\n", err)
			fmt.Println("Continuing with Cloudflare deploy...")
		}
	}

	// 5. Deploy via wrangler
	fmt.Printf("Deploying %s to %s.%s...\n", cfg.Name, cfg.Deploy.Subdomain, cfg.Deploy.Domain)
	if err := deploy.DeployWithWrangler(distDir, cfg.Deploy.Subdomain); err != nil {
		return err
	}

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
