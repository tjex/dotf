package config

import (
	"os"
	"slices"
	"testing"
)

// test config.extractSubmodules()
func TestSubmoduleExtraction(t *testing.T) {
	testConf := "./testdata/test-config"
	cfg.Worktree = "/Users/foo"
	file, err := os.Open(testConf)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"/Users/foo/.config/asciinema",
		"/Users/foo/.config/aerc",
		"/Users/foo/.config/goimap",
		// this is transformed from a relative filepath in test-config
		"/Users/foo/.local/share/bkmr",
	}
	have := extractSubmodulePaths(file)

	if slices.Compare(want, have) != 0 {
		t.Fatal("\nwant:", want, "\nbut have:", have)
	}
}
