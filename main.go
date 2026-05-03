package main

import (
	"embed"

	"github.com/mrmackaniel/launchpad/cmd"
	tmpl "github.com/mrmackaniel/launchpad/internal/template"
)

//go:embed templates
var builtinTemplates embed.FS

func main() {
	tmpl.BuiltinTemplates = builtinTemplates
	cmd.Execute()
}
