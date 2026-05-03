package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/mrmackaniel/launchpad/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Scaffold a showcase.config.json interactively",
	RunE:  runInit,
}

var initOutput string

func init() {
	initCmd.Flags().StringVarP(&initOutput, "output", "o", "showcase.config.json", "output file path")
}

func runInit(cmd *cobra.Command, args []string) error {
	if _, err := os.Stat(initOutput); err == nil {
		var overwrite bool
		if err := huh.NewConfirm().
			Title(initOutput + " already exists. Overwrite?").
			Value(&overwrite).
			Run(); err != nil {
			return err
		}
		if !overwrite {
			fmt.Println("Aborted.")
			return nil
		}
	}

	var (
		name        string
		tagline     string
		description string
		featuresRaw string
		quickStart  string
		techRaw     string
		github      string
		live        string
		videoURL    string
		subdomain   string
		domain      string
		template    string
	)

	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Project name").Value(&name).Validate(required("name")),
			huh.NewInput().Title("Tagline").Description("One-liner shown under the title").Value(&tagline).Validate(required("tagline")),
			huh.NewText().Title("Description").Description("A paragraph introducing the project").Value(&description).Validate(required("description")),
		),
		huh.NewGroup(
			huh.NewText().Title("Features").Description("One feature per line").Value(&featuresRaw),
			huh.NewInput().Title("Quick start command").Description("e.g. go install github.com/you/project@latest").Value(&quickStart),
			huh.NewInput().Title("Tech stack").Description("Comma-separated, e.g. Go, Docker, Postgres").Value(&techRaw),
		),
		huh.NewGroup(
			huh.NewInput().Title("GitHub URL").Value(&github),
			huh.NewInput().Title("Live site URL").Description("Optional — leave blank if not deployed").Value(&live),
			huh.NewInput().Title("Demo video URL").Description("YouTube, Vimeo, or .mp4 path — optional").Value(&videoURL),
		),
		huh.NewGroup(
			huh.NewSelect[string]().Title("Template").Options(
				huh.NewOption("minimal — clean, typography-focused", "minimal"),
				huh.NewOption("bold — dark hero, high contrast", "bold"),
				huh.NewOption("developer — terminal/code aesthetic", "developer"),
			).Value(&template),
			huh.NewInput().Title("Subdomain").Description("e.g. myproject → myproject.yourdomain.com").Value(&subdomain).Validate(required("subdomain")),
			huh.NewInput().Title("Domain").Description("e.g. yourdomain.com").Value(&domain).Validate(required("domain")),
		),
	).Run(); err != nil {
		return err
	}

	features := splitLines(featuresRaw)
	techStack := splitCommas(techRaw)

	cfg := config.Config{
		Name:        name,
		Tagline:     tagline,
		Description: description,
		Features:    features,
		QuickStart:  quickStart,
		TechStack:   techStack,
		Media:       config.Media{Video: videoURL},
		Links:       config.Links{GitHub: github, Live: live},
		Template:    template,
		Deploy: config.Deploy{
			Subdomain: subdomain,
			Domain:    domain,
			Provider:  "cloudflare",
		},
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(initOutput, data, 0644); err != nil {
		return err
	}

	fmt.Printf("\nCreated %s\n", initOutput)
	fmt.Printf("Run `launchpad preview` to see it locally.\n")
	return nil
}

func required(field string) func(string) error {
	return func(v string) error {
		if strings.TrimSpace(v) == "" {
			return fmt.Errorf("%s is required", field)
		}
		return nil
	}
}

func splitLines(s string) []string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		if t := strings.TrimSpace(line); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func splitCommas(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		if t := strings.TrimSpace(part); t != "" {
			out = append(out, t)
		}
	}
	return out
}
