package config

import (
	"bufio"
	"os"
	"slices"
	"strings"
	"testing"
)

func TestSubmodules(t *testing.T) {
	testConf := "../../fixtures/test-config"
	var configLines, submodulePaths []string
	file, err := os.Open(testConf)
	if err != nil {
		t.Fatal("error opening bare repository config:", err)
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		configLines = append(configLines, sc.Text())
	}
	for _, line := range configLines {
		match := submLineRe.FindString(line)
		submodulePath := submPathRe.FindString(match)
		if submodulePath != "" {
			submodulePath = strings.ReplaceAll(submodulePath, `"`, "")
			submodulePaths = append(submodulePaths, submodulePath)
		}
	}
	want := []string{
		"/Users/foo/.config/asciinema",
		"/Users/foo/.config/aerc",
		"/Users/foo/.config/goimap",
	}
	have := submodulePaths
	if slices.Compare(want, have) != 0 {
		t.Fatal("\nwant:", want, "\nhave:", have)
	}
}
