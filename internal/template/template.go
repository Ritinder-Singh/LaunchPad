package template

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"strings"

	"github.com/mrmackaniel/launchpad/internal/config"
)

// BuiltinTemplates is set from main.go via go:embed so the embed directive
// can reference the templates/ dir from the module root.
var BuiltinTemplates fs.ReadFileFS

type templateData struct {
	*config.Config
	VideoEmbedURL string
	VideoType     string
	Screenshots   []string // relative paths from dist root, e.g. "screenshots/foo.png"
}

func Render(cfg *config.Config, screenshots []string) ([]byte, error) {
	var raw []byte
	var err error

	if strings.HasPrefix(cfg.Template, ".") || strings.HasPrefix(cfg.Template, "/") {
		raw, err = os.ReadFile(cfg.Template)
		if err != nil {
			return nil, fmt.Errorf("loading custom template: %w", err)
		}
	} else {
		path := fmt.Sprintf("templates/%s/index.html", cfg.Template)
		raw, err = BuiltinTemplates.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("unknown template %q: %w", cfg.Template, err)
		}
	}

	tmpl, err := template.New("showcase").Parse(string(raw))
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	embedURL, videoType := resolveVideo(cfg.Media.Video)
	data := templateData{
		Config:        cfg,
		VideoEmbedURL: embedURL,
		VideoType:     videoType,
		Screenshots:   screenshots,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering template: %w", err)
	}

	return buf.Bytes(), nil
}

func resolveVideo(v string) (embedURL, videoType string) {
	if v == "" {
		return "", ""
	}
	if strings.HasSuffix(v, ".mp4") {
		return v, "mp4"
	}
	if strings.Contains(v, "youtube.com/watch?v=") {
		id := v[strings.Index(v, "v=")+2:]
		if i := strings.IndexAny(id, "&?"); i != -1 {
			id = id[:i]
		}
		return "https://www.youtube.com/embed/" + id, "youtube"
	}
	if strings.Contains(v, "youtu.be/") {
		id := v[strings.LastIndex(v, "/")+1:]
		if i := strings.IndexAny(id, "?&"); i != -1 {
			id = id[:i]
		}
		return "https://www.youtube.com/embed/" + id, "youtube"
	}
	if strings.Contains(v, "vimeo.com/") {
		id := v[strings.LastIndex(v, "/")+1:]
		if i := strings.IndexAny(id, "?&"); i != -1 {
			id = id[:i]
		}
		return "https://player.vimeo.com/video/" + id, "vimeo"
	}
	return "", ""
}
