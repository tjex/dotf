package config

import (
	"os"
	"slices"
	"testing"
)

// test config.extractSubmodules()
func TestSubmoduleExtraction(t *testing.T) {
	testConf := "./testdata/test-config"
	file, err := os.Open(testConf)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"/Users/foo/.config/asciinema",
		"/Users/foo/.config/aerc",
		"/Users/foo/.config/goimap",
	}
	have := extractSubmodules(file)

	if slices.Compare(want, have) != 0 {
		t.Fatal("\nwant:", want, "\nbut have:", have)
	}
}
