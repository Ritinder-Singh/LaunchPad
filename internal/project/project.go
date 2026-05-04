package project

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Type string

const (
	TypeNode   Type = "node"
	TypeGo     Type = "go"
	TypePython Type = "python"
	TypeStatic Type = "static"
)

type Project struct {
	Type         Type
	Dir          string
	BuildCmd     string
	StartCmd     string
	Port         int
	proc         *exec.Cmd
}

// Detect inspects dir and returns a Project with sensible defaults.
func Detect(dir string) *Project {
	p := &Project{Dir: dir}

	if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
		p.Type = TypeNode
		p.BuildCmd = "npm install && npm run build"
		p.StartCmd = "npm start"
		p.Port = 3000
	} else if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
		p.Type = TypeGo
		p.BuildCmd = "go build -o .showcase/.bin/app ."
		p.StartCmd = "./.showcase/.bin/app"
		p.Port = 8080
	} else if _, err := os.Stat(filepath.Join(dir, "requirements.txt")); err == nil {
		p.Type = TypePython
		p.BuildCmd = "pip install -r requirements.txt"
		p.StartCmd = "python app.py"
		p.Port = 5000
	} else {
		p.Type = TypeStatic
		p.BuildCmd = ""
		p.StartCmd = fmt.Sprintf("python -m http.server 8000")
		p.Port = 8000
	}

	return p
}

// Apply overrides from config values (zero values are ignored).
func (p *Project) Apply(buildCmd, startCmd string, port int) {
	if buildCmd != "" {
		p.BuildCmd = buildCmd
	}
	if startCmd != "" {
		p.StartCmd = startCmd
	}
	if port != 0 {
		p.Port = port
	}
}

// Build runs the build command, if any.
func (p *Project) Build() error {
	if p.BuildCmd == "" {
		return nil
	}
	fmt.Printf("Building (%s)...\n", p.BuildCmd)
	cmd := exec.Command("sh", "-c", p.BuildCmd)
	cmd.Dir = p.Dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Start launches the project's start command in the background.
func (p *Project) Start() error {
	if p.StartCmd == "" {
		return fmt.Errorf("no start command configured")
	}
	fmt.Printf("Starting project on port %d...\n", p.Port)
	p.proc = exec.Command("sh", "-c", p.StartCmd)
	p.proc.Dir = p.Dir
	p.proc.Stdout = os.Stdout
	p.proc.Stderr = os.Stderr
	return p.proc.Start()
}

// Stop kills the running project process.
func (p *Project) Stop() {
	if p.proc != nil && p.proc.Process != nil {
		p.proc.Process.Kill()
		p.proc.Wait()
		p.proc = nil
	}
}

// WaitReady polls localhost:port until reachable or timeout.
func (p *Project) WaitReady(timeout time.Duration) error {
	addr := fmt.Sprintf("localhost:%d", p.Port)
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, time.Second)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("project did not start on port %d within %s", p.Port, timeout)
}
