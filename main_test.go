package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestSyncReportsRepositoryErrors(t *testing.T) {
	dir := t.TempDir()
	binary := filepath.Join(dir, "dotf")
	build := exec.Command("go", "build", "-buildvcs=false", "-o", binary, ".")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, output)
	}

	bare := filepath.Join(dir, "missing-bare")
	first := filepath.Join(dir, "missing-first")
	second := filepath.Join(dir, "missing-second")
	for _, path := range []string{bare, first, second} {
		if err := os.Mkdir(path, 0755); err != nil {
			t.Fatal(err)
		}
	}
	configDir := filepath.Join(dir, "config", "dotf")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(`git-dir = %q
worktree = %q
[modules.default]
paths = [%q, %q]
`, bare, dir, first, second)
	if err := os.WriteFile(filepath.Join(configDir, "config.toml"), []byte(config), 0600); err != nil {
		t.Fatal(err)
	}

	for _, tt := range []struct {
		name  string
		args  []string
		repos []string
	}{
		{"bare", []string{"sync", "-b", "-q"}, []string{bare}},
		{"modules", []string{"sync", "-m", "-q"}, []string{first, second}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			cmd.Env = append(os.Environ(), "XDG_CONFIG_HOME="+filepath.Join(dir, "config"))
			output, _ := cmd.CombinedOutput()
			text := string(output)
			for _, repo := range tt.repos {
				if !strings.Contains(text, repo) {
					t.Errorf("missing repository %q in output:\n%s", repo, text)
				}
				if !strings.Contains(text, repo+") fatal:") {
					t.Errorf("repository %q not associated with Git error:\n%s", repo, text)
				}
			}
			if !strings.Contains(text, "fatal:") {
				t.Errorf("missing Git error in output:\n%s", text)
			}
			if strings.Contains(text, "Syncing ") {
				t.Errorf("quiet mode printed progress:\n%s", text)
			}
		})
	}

	t.Run("all", func(t *testing.T) {
		remote := filepath.Join(dir, "remote.git")
		worktree := filepath.Join(dir, "worktree")
		for _, args := range [][]string{
			{"init", "--bare", remote},
			{"clone", remote, worktree},
			{"-C", worktree, "-c", "user.name=Test", "-c", "user.email=test@example.com", "commit", "--allow-empty", "-m", "initial"},
			{"-C", worktree, "push", "-u", "origin", "HEAD"},
		} {
			cmd := exec.Command("git", args...)
			cmd.Env = append(os.Environ(), "HOME="+dir, "GIT_CONFIG_NOSYSTEM=1")
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("git %v: %v\n%s", args, err, output)
			}
		}
		gitDir := filepath.Join(worktree, ".git")
		config := fmt.Sprintf(`git-dir = %q
worktree = %q
[modules.default]
paths = [%q, %q]
`, gitDir, worktree, first, second)
		if err := os.WriteFile(filepath.Join(configDir, "config.toml"), []byte(config), 0600); err != nil {
			t.Fatal(err)
		}
		cmd := exec.Command(binary, "sync", "-a", "-q")
		cmd.Env = append(os.Environ(), "HOME="+dir, "GIT_CONFIG_NOSYSTEM=1", "XDG_CONFIG_HOME="+filepath.Join(dir, "config"))
		output, _ := cmd.CombinedOutput()
		for _, repo := range []string{first, second} {
			if !strings.Contains(string(output), "("+repo+") fatal:") {
				t.Errorf("missing error for %s:\n%s", repo, output)
			}
		}
		if strings.Contains(string(output), "Syncing ") {
			t.Errorf("quiet mode printed progress:\n%s", output)
		}
	})
}
