package capture

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Screenshots returns paths to screenshot images, capturing them via Puppeteer
// if the screenshotDir has no existing images and autoCapture is enabled.
func Screenshots(screenshotDir string, autoCapture bool, start func() error, waitReady func(time.Duration) error, stop func()) ([]string, error) {
	// Check for existing screenshots first
	existing, _ := findImages(screenshotDir)
	if len(existing) > 0 {
		fmt.Printf("Using %d existing screenshot(s) from %s\n", len(existing), screenshotDir)
		return existing, nil
	}

	if !autoCapture {
		fmt.Printf("No screenshots found in %s and autoCapture is disabled.\n", screenshotDir)
		fmt.Println("Add images to that folder and re-run, or set autoCapture: true in config.")
		return nil, nil
	}

	fmt.Println("No screenshots found — capturing via Puppeteer...")

	if err := start(); err != nil {
		return nil, fmt.Errorf("starting project for screenshot: %w", err)
	}
	defer stop()

	fmt.Println("Waiting for project to be ready...")
	if err := waitReady(30 * time.Second); err != nil {
		return nil, fmt.Errorf("project did not start: %w\nTip: add screenshots manually to %s to skip auto-capture", err, screenshotDir)
	}

	if err := os.MkdirAll(screenshotDir, 0755); err != nil {
		return nil, err
	}

	if err := runPuppeteer(screenshotDir); err != nil {
		fmt.Printf("Warning: auto-capture failed: %v\n", err)
		fmt.Printf("Tip: add screenshots manually to %s and re-run.\n", screenshotDir)
		return nil, nil
	}

	return findImages(screenshotDir)
}

func runPuppeteer(screenshotDir string) error {
	// Write capture script to a temp file under .showcase/
	scriptPath := filepath.Join(".showcase", ".capture.js")

	absDir, err := filepath.Abs(screenshotDir)
	if err != nil {
		return err
	}

	script := fmt.Sprintf(`
const puppeteer = require('puppeteer');
const path = require('path');

(async () => {
  const browser = await puppeteer.launch({
    args: ['--no-sandbox', '--disable-setuid-sandbox'],
  });
  const page = await browser.newPage();
  await page.setViewport({ width: 1280, height: 800 });
  try {
    await page.goto(process.env.CAPTURE_URL || 'http://localhost:3000', {
      waitUntil: 'networkidle2',
      timeout: 20000,
    });
    // Wait a moment for animations/lazy images
    await new Promise(r => setTimeout(r, 1500));
    await page.screenshot({
      path: path.join(%q, 'screenshot.png'),
      fullPage: false,
    });
    console.log('Screenshot saved to %s/screenshot.png');
  } catch (e) {
    console.error('Screenshot error:', e.message);
    process.exit(1);
  } finally {
    await browser.close();
  }
})();
`, absDir, screenshotDir)

	if err := os.WriteFile(scriptPath, []byte(script), 0644); err != nil {
		return fmt.Errorf("writing capture script: %w", err)
	}
	defer os.Remove(scriptPath)

	// Run via npx so puppeteer is auto-installed if missing
	cmd := exec.Command("npx", "--yes", "--package=puppeteer", "node", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func findImages(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		switch filepath.Ext(e.Name()) {
		case ".png", ".jpg", ".jpeg", ".webp", ".gif":
			out = append(out, filepath.Join(dir, e.Name()))
		}
	}
	return out, nil
}
