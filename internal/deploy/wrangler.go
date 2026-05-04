package deploy

import (
	"fmt"
	"os"
	"os/exec"
)

// DeployWithWrangler deploys distDir to Cloudflare Pages via npx wrangler.
// Requires CLOUDFLARE_API_TOKEN and CLOUDFLARE_ACCOUNT_ID in the environment.
func DeployWithWrangler(distDir, projectName string) error {
	fmt.Printf("Deploying %s via wrangler...\n", projectName)

	cmd := exec.Command("npx", "--yes", "wrangler", "pages", "deploy", distDir,
		"--project-name", projectName,
		"--commit-dirty=true",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wrangler deploy: %w", err)
	}
	return nil
}
