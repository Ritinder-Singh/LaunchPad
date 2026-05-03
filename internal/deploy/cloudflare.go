package deploy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mrmackaniel/launchpad/internal/config"
)

type Deployer interface {
	Deploy(cfg *config.Config, distDir string) (string, error)
}

func New(provider string) (Deployer, error) {
	switch provider {
	case "cloudflare", "":
		token := os.Getenv("CLOUDFLARE_API_TOKEN")
		accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
		if token == "" || accountID == "" {
			return nil, fmt.Errorf("CLOUDFLARE_API_TOKEN and CLOUDFLARE_ACCOUNT_ID must be set")
		}
		return &CloudflareDeployer{token: token, accountID: accountID}, nil
	default:
		return nil, fmt.Errorf("unsupported provider %q — see DEPLOYERS.md", provider)
	}
}

type CloudflareDeployer struct {
	token     string
	accountID string
}

func (d *CloudflareDeployer) Deploy(cfg *config.Config, distDir string) (string, error) {
	projectName := cfg.Deploy.Subdomain

	if err := d.ensureProject(projectName, cfg.Deploy.Domain); err != nil {
		return "", fmt.Errorf("ensuring Cloudflare Pages project: %w", err)
	}

	url, err := d.uploadFiles(projectName, distDir)
	if err != nil {
		return "", fmt.Errorf("uploading files: %w", err)
	}

	return url, nil
}

func (d *CloudflareDeployer) ensureProject(name, domain string) error {
	body, _ := json.Marshal(map[string]any{
		"name":              name,
		"production_branch": "main",
	})

	req, _ := http.NewRequest("POST",
		fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/pages/projects", d.accountID),
		bytes.NewReader(body),
	)
	d.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 409 means project already exists — that's fine
	if resp.StatusCode != 200 && resp.StatusCode != 409 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, b)
	}
	return nil
}

func (d *CloudflareDeployer) uploadFiles(projectName, distDir string) (string, error) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)

	err := filepath.WalkDir(distDir, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(distDir, path)
		fw, err := w.CreateFormFile(rel, rel)
		if err != nil {
			return err
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(fw, f)
		return err
	})
	if err != nil {
		return "", err
	}
	w.Close()

	req, _ := http.NewRequest("POST",
		fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/pages/projects/%s/deployments",
			d.accountID, projectName),
		&body,
	)
	d.setHeaders(req)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Result struct {
			URL string `json:"url"`
		} `json:"result"`
		Success bool `json:"success"`
		Errors  []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if !result.Success {
		if len(result.Errors) > 0 {
			return "", fmt.Errorf("Cloudflare API error: %s", result.Errors[0].Message)
		}
		return "", fmt.Errorf("Cloudflare API returned failure")
	}

	return result.Result.URL, nil
}

func (d *CloudflareDeployer) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+d.token)
}
