# Launchpad

Deploy a project showcase page to your own subdomain from a single JSON config.

## What It Is

`launchpad` is a Go CLI tool that reads a `showcase.config.json`, generates a static project showcase page from a built-in (or custom) template, and deploys it automatically to a custom subdomain via Cloudflare Pages.

The tool itself is showcased as a project — its own page was generated and deployed using `launchpad`.

## How It Works

1. Create a `showcase.config.json` for your project
2. Run `launchpad deploy`
3. Your project page is live at `projectname.yourdomain.com`

## Commands

```
launchpad init              # scaffold showcase.config.json interactively
launchpad preview           # serve the generated page locally
launchpad deploy            # build + deploy to Cloudflare Pages
launchpad templates list    # list available built-in templates
```

## Config Format

```json
{
  "name": "My Project",
  "tagline": "Short one-liner",
  "description": "Longer intro paragraph...",
  "techStack": ["Go", "Astro", "Cloudflare Pages"],
  "images": ["./screenshots/hero.png", "./screenshots/demo.png"],
  "links": {
    "github": "https://github.com/you/project",
    "live": "https://project.yourdomain.com"
  },
  "template": "minimal",
  "deploy": {
    "subdomain": "myproject",
    "domain": "yourdomain.com",
    "provider": "cloudflare"
  }
}
```

`links.live` is optional — omit it if the project isn't deployed anywhere.

`template` accepts a built-in name (`minimal`, `bold`, `developer`) or a path to your own HTML file (`./my-template.html`).

## Built-in Templates

| Name | Description |
|---|---|
| `minimal` | Clean, whitespace-heavy, typography-focused |
| `bold` | Dark background, large hero, accent color |
| `developer` | Terminal/code aesthetic, monospace font |

## Stack

- **CLI:** Go — single binary, no runtime dependencies
- **Templates:** Pre-built HTML/CSS embedded in the binary via `go:embed`
- **Deployment:** Cloudflare Pages (primary) — see [DEPLOYERS.md](./DEPLOYERS.md) for planned providers

## Deployment

Requires a Cloudflare API token with Pages and DNS permissions:

```
export CLOUDFLARE_API_TOKEN=your_token
export CLOUDFLARE_ACCOUNT_ID=your_account_id

launchpad deploy
```
