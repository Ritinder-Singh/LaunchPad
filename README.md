# Launchpad

A CLI tool that turns a JSON config into a live project showcase page on your own subdomain in under 60 seconds.

## What It Is

`launchpad` is a deployment tool that generates static project showcase websites from a simple config file and deploys them automatically to a custom subdomain. Each project gets its own dedicated page at `projectname.yourdomain.com`.

The tool itself is showcased as a project — its own page was generated and deployed using `launchpad`.

## How It Works

1. Create a `showcase.config.json` for your project
2. Run `npx showcase deploy`
3. Your project page is live at `projectname.yourdomain.com`

Each page includes a hero section, project introduction, screenshots, tech stack, and links — all driven by the config file.

## Stack

- **Template:** Astro — outputs pure static HTML, no JS overhead
- **CLI:** Node.js — distributable via `npx`
- **Deployment:** Cloudflare Pages — free tier, fast subdomain setup via API
- **Config:** `showcase.config.json` — one file per project

## Config Format

```json
{
  "name": "My Project",
  "tagline": "A short description of what this does",
  "description": "A longer introduction to the project...",
  "techStack": ["TypeScript", "Astro", "Cloudflare Pages"],
  "images": ["./screenshots/hero.png", "./screenshots/demo.png"],
  "links": {
    "github": "https://github.com/you/project",
    "live": "https://project.yourdomain.com"  // optional — omit if not deployed
  }
}
```

## Deployment

Each project deploys to its own subdomain via Cloudflare Pages. Subdomains are configured automatically through the Cloudflare API.

```
launchpad deploy --config ./showcase.config.json --subdomain myproject
```
