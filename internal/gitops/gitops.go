package gitops

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const showcaseBranch = "showcase"
const worktreePath = ".showcase/.worktree"

// PushShowcase commits the contents of distDir to the `showcase` branch and pushes.
func PushShowcase(distDir string) error {
	// Ensure showcase branch exists locally
	if err := ensureShowcaseBranch(); err != nil {
		return fmt.Errorf("preparing showcase branch: %w", err)
	}

	// Clean up any stale worktree
	os.RemoveAll(worktreePath)
	gitRun("worktree", "prune")

	// Add a worktree for the showcase branch
	if err := gitRun("worktree", "add", worktreePath, showcaseBranch); err != nil {
		return fmt.Errorf("creating git worktree: %w", err)
	}
	defer func() {
		gitRun("worktree", "remove", "--force", worktreePath)
		os.RemoveAll(worktreePath)
	}()

	// Clear worktree content (preserve .git file)
	if err := clearWorktree(worktreePath); err != nil {
		return fmt.Errorf("clearing worktree: %w", err)
	}

	// Copy dist files into the worktree
	if err := copyDir(distDir, worktreePath); err != nil {
		return fmt.Errorf("copying showcase files: %w", err)
	}

	// Commit
	gitInDir(worktreePath, "add", "-A")
	if err := gitInDir(worktreePath, "commit", "-m", "chore: update showcase [skip ci]", "--allow-empty"); err != nil {
		return fmt.Errorf("committing showcase: %w", err)
	}

	// Push
	fmt.Println("Pushing showcase branch to origin...")
	if err := gitInDir(worktreePath, "push", "origin", showcaseBranch, "--force"); err != nil {
		return fmt.Errorf("pushing showcase: %w", err)
	}

	fmt.Printf("Showcase pushed to branch %q\n", showcaseBranch)
	return nil
}

func ensureShowcaseBranch() error {
	// Check if branch already exists locally
	err := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+showcaseBranch).Run()
	if err == nil {
		return nil // already exists
	}

	// Check if it exists on remote and fetch it
	remoteCheck := exec.Command("git", "ls-remote", "--exit-code", "--heads", "origin", showcaseBranch)
	if remoteCheck.Run() == nil {
		return gitRun("fetch", "origin", showcaseBranch+":"+showcaseBranch)
	}

	// Create a new orphan branch without switching the working tree
	fmt.Printf("Creating new %q branch...\n", showcaseBranch)

	// Get an empty tree hash
	out, err := exec.Command("git", "hash-object", "-t", "tree", "/dev/null").Output()
	if err != nil {
		// Fallback: create branch from current HEAD (simpler but carries project history)
		return gitRun("branch", showcaseBranch)
	}
	emptyTree := strings.TrimSpace(string(out))

	// Create empty commit
	commitOut, err := exec.Command("git", "commit-tree", emptyTree, "-m", "init showcase branch").Output()
	if err != nil {
		return gitRun("branch", showcaseBranch) // fallback
	}
	commitHash := strings.TrimSpace(string(commitOut))

	return gitRun("branch", showcaseBranch, commitHash)
}

func clearWorktree(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.Name() == ".git" {
			continue
		}
		os.RemoveAll(filepath.Join(dir, e.Name()))
	}
	return nil
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func gitRun(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func gitInDir(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
