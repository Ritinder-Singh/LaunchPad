# TODO

## Documentation
- [ ] Add README badges (CI status, latest release, open issues, license)
- [ ] Add Homebrew formula / `go install` instructions to README

## Templates
- [ ] Polish `bold` template design
- [ ] Polish `developer` template design
- [ ] Set up `launchpad-templates` GitHub repo (see TEMPLATE_REPO.md)
- [ ] Wire CLI to fetch templates from remote repo with local cache (`~/.launchpad/cache/`)

## CLI
- [ ] Polish `launchpad init` interactive prompts (validation, better UX)
- [ ] Test `launchpad deploy` with real Cloudflare credentials

## Infrastructure
- [ ] Set branch protection rules on GitHub:
  - `main`: require PR, no direct push, require CI to pass
  - `develop`: no direct push to main from here
- [ ] Verify release workflow fires correctly on tag push

## Distribution
- [ ] Publish to GitHub Releases with a `v0.1.0` tag
- [ ] Homebrew formula (after first release)
