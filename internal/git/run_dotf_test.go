package git_test

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"git.sr.ht/~tjex/dotf/internal/config"
	"git.sr.ht/~tjex/dotf/internal/git"
)

func TestRunDotfShowsInteractiveStderr(t *testing.T) {
	if os.Getenv("DOTF_TEST_GIT_HELPER") == "1" {
		config.UserConfig().GitDir = "/example/repo"
		_ = git.RunDotf([]string{"push"})
		return
	}

	dir := t.TempDir()
	script := "#!/bin/sh\nprintf 'Password: ' >&2\nread password\nprintf '\\nfatal: authentication failed\\n' >&2\nexit 1\n"
	if err := os.WriteFile(filepath.Join(dir, "git"), []byte(script), 0755); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestRunDotfShowsInteractiveStderr$")
	cmd.Env = append(os.Environ(), "DOTF_TEST_GIT_HELPER=1", "PATH="+dir)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	prompt := make(chan string, 1)
	go func() {
		text, _ := bufio.NewReader(stderr).ReadString(':')
		prompt <- text
	}()
	select {
	case text := <-prompt:
		if !strings.Contains(text, "/example/repo") || !strings.Contains(text, "Password:") {
			t.Errorf("prompt lacks repository context: %q", text)
		}
		_, _ = stdin.Write([]byte("secret\n"))
		_ = stdin.Close()
	case <-time.After(2 * time.Second):
		_ = cmd.Process.Kill()
		t.Error("Git prompt was not visible before input")
	}
	_ = cmd.Wait()
}

func TestSyncStateShowsInteractiveStderr(t *testing.T) {
	if os.Getenv("DOTF_TEST_FETCH_HELPER") == "1" {
		_, _, _ = git.SyncState("/example/repo")
		return
	}

	dir := t.TempDir()
	script := "#!/bin/sh\nprintf 'Password: ' >&2\nread password\nprintf '\\nfatal: authentication failed\\n' >&2\nexit 1\n"
	if err := os.WriteFile(filepath.Join(dir, "git"), []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestSyncStateShowsInteractiveStderr$")
	cmd.Env = append(os.Environ(), "DOTF_TEST_FETCH_HELPER=1", "PATH="+dir)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	prompt := make(chan string, 1)
	go func() {
		text, _ := bufio.NewReader(stderr).ReadString(':')
		prompt <- text
	}()
	select {
	case text := <-prompt:
		if !strings.Contains(text, "/example/repo") || !strings.Contains(text, "Password:") {
			t.Errorf("prompt lacks repository context: %q", text)
		}
		_, _ = stdin.Write([]byte("secret\n"))
		_ = stdin.Close()
	case <-time.After(2 * time.Second):
		_ = cmd.Process.Kill()
		t.Error("Git prompt was not visible before input")
	}
	_ = cmd.Wait()
}
