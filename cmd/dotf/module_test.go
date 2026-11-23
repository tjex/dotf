package dotf

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"git.sr.ht/~tjex/dotf/internal/config"
	"git.sr.ht/~tjex/dotf/internal/printer"
)

func TestPullReturnsEveryRepositoryFailure(t *testing.T) {
	dir := t.TempDir()
	first, second := filepath.Join(dir, "first"), filepath.Join(dir, "second")
	for _, path := range []string{first, second} {
		if err := os.Mkdir(path, 0755); err != nil {
			t.Fatal(err)
		}
	}

	previous := *cfg
	t.Cleanup(func() { *cfg = previous })
	cfg.Modules = map[string]config.ModuleConfig{
		"default": {Paths: []config.Module{{Path: first}, {Path: second}}},
	}
	cfg.Concurrency.BatchSize = 2
	modules := &Modules{Printer: printer.NewPrinter(true)}
	err := modules.Pull()
	if err == nil {
		t.Fatal("expected repository failures")
	}
	for _, repo := range []string{first, second} {
		if !strings.Contains(err.Error(), "("+repo+") fatal:") {
			t.Errorf("missing failure for %s: %v", repo, err)
		}
	}
}
