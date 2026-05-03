# launchpad-templates

> Notes for setting up the companion template repository.

## What This Repo Is

`launchpad-templates` is a separate GitHub repository that holds the design templates used by the Launchpad CLI. Each template is a standalone folder containing an `index.html` (and optionally CSS/JS/assets) with Go `html/template` placeholders.

Separating templates from the CLI means:
- Templates can be updated without releasing a new CLI version
- The community can contribute new designs via pull requests

## Repository to Create

**`github.com/[username]/launchpad-templates`**

## Directory Structure

```
launchpad-templates/
├── minimal/
│   └── index.html
├── bold/
│   └── index.html
├── developer/
│   └── index.html
└── CONTRIBUTING.md
```

## Template Variables

Every template receives the following data from `showcase.config.json`:

| Variable | Type | Notes |
|---|---|---|
| `{{.Name}}` | string | Project name |
| `{{.Tagline}}` | string | Short one-liner |
| `{{.Description}}` | string | Longer intro paragraph |
| `{{.TechStack}}` | []string | List of tech/tools used |
| `{{.Images}}` | []string | Paths or URLs to screenshots |
| `{{.Links.GitHub}}` | string | GitHub repo URL |
| `{{.Links.Live}}` | string | Live URL — may be empty, always check with `{{if .Links.Live}}` |

## Contributing a Template

1. Fork the repo
2. Create a new folder with your template name (lowercase, no spaces)
3. Add an `index.html` using the variables above
4. Open a PR with a screenshot of what it looks like

## How the CLI Uses This Repo

The CLI fetches templates from this repo at deploy time and caches them locally at `~/.launchpad/cache/`. Users can pin a version:

```json
"template": "minimal@v1.2"
```

## TODO

- [ ] Set up the GitHub repo
- [ ] Add CONTRIBUTING.md with screenshot requirements
- [ ] Set up versioning/tagging convention for template releases
- [ ] Wire up CLI to fetch from this repo (currently uses embedded templates)
