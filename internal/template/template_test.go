package template_test

import (
	"embed"
	"io/fs"
	"strings"
	"testing"

	"github.com/mrmackaniel/launchpad/internal/config"
	tmpl "github.com/mrmackaniel/launchpad/internal/template"
)

//go:embed testdata
var testTemplates embed.FS

func subFS(t *testing.T) fs.ReadFileFS {
	t.Helper()
	sub, err := fs.Sub(testTemplates, "testdata")
	if err != nil {
		t.Fatal(err)
	}
	return sub.(fs.ReadFileFS)
}

func TestRenderBuiltin(t *testing.T) {
	tmpl.BuiltinTemplates = subFS(t)

	cfg := &config.Config{
		Name:        "Test Project",
		Tagline:     "A test tagline",
		Description: "A test description.",
		Features:    []string{"Fast", "Simple", "Reliable"},
		QuickStart:  "go install github.com/test/project@latest",
		TechStack:   []string{"Go", "HTML"},
		Links:       config.Links{GitHub: "https://github.com/test/project"},
		Template:    "minimal",
	}

	out, err := tmpl.Render(cfg)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	for _, want := range []string{
		"Test Project", "A test tagline", "A test description.",
		"Fast", "Simple", "Reliable",
		"go install github.com/test/project@latest",
		"Go", "HTML",
	} {
		if !strings.Contains(string(out), want) {
			t.Errorf("output missing %q", want)
		}
	}

	if strings.Contains(string(out), "Live site") {
		t.Error("optional live link should not render when empty")
	}
}

func TestResolveVideoYouTube(t *testing.T) {
	tmpl.BuiltinTemplates = subFS(t)

	cfg := &config.Config{
		Name:     "Video Test",
		Template: "minimal",
		Media:    config.Media{Video: "https://www.youtube.com/watch?v=dQw4w9WgXcQ"},
	}

	out, err := tmpl.Render(cfg)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(string(out), "youtube.com/embed/dQw4w9WgXcQ") {
		t.Error("YouTube URL not converted to embed URL")
	}
}

func TestResolveVideoMP4(t *testing.T) {
	tmpl.BuiltinTemplates = subFS(t)

	cfg := &config.Config{
		Name:     "MP4 Test",
		Template: "minimal",
		Media:    config.Media{Video: "./demo.mp4"},
	}

	out, err := tmpl.Render(cfg)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(string(out), "<video") {
		t.Error("mp4 should render a <video> tag")
	}
}
