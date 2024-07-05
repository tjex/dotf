package config

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
)

var (
	buf        bytes.Buffer
	logger     = log.New(&buf, "logger: ", log.Lshortfile)
	cfg        Config
	submLineRe = regexp.MustCompile(`^\[submodule ".+?"\]`)
	submPathRe = regexp.MustCompile(`"(.+?)"`)
)

type Config struct {
	Worktree           string
	GitDir             string `toml:"git-dir"`
	Origin             string
	Mirror             string
	RepoFlags          []string // Flags for bare repo dir and worktree
	BatchCommitMessage string   `toml:"batch-commit-message"`
}

// returns pointer to user config struct
func UserConfig() *Config {
	c := &cfg
	return c
}

// reads user configuration from a .toml file
func ReadUserConfig() {
	path := configDir()
	confFile := filepath.Join(path, "config.toml")
	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		fmt.Println("No config file at \"XDG_CONFIG_HOME/dotf/config.toml\", does it exist?")
		os.Exit(1)
	}
	file, err := os.ReadFile(confFile)
	if err != nil {
		logger.Print(err)
	}
	parseConfig(file)

}

// parse a .toml file for usage.
func parseConfig(config []byte) {
	err := toml.Unmarshal(config, &cfg)
	if err != nil {
		logger.Print(err)
	}
	buildGitRepoFlags(&cfg)
}

// build the bare repo git argument array
func buildGitRepoFlags(cfg *Config) {
	var flags []string
	flags = append(flags, "--git-dir", cfg.GitDir, "--work-tree", cfg.Worktree)
	cfg.RepoFlags = flags
}

func configDir() string {
	path, ok := os.LookupEnv("XDG_CONFIG_HOME")
	if !ok {
		home, ok := os.LookupEnv("HOME")
		if !ok {
			home = "~/"
		}
		path = filepath.Join(home, ".config")
	}
	return filepath.Join(path, "dotf")
}

// Retrieve path to bare repo git config
func BareRepoConfig() string {
	gitConf := filepath.Join(cfg.GitDir, "config")
	return gitConf

}

// Parse the provided git config file and return its submodule paths.
func extractSubmodulePaths(file *os.File) []string {
	var configLines, submodulePaths []string
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		configLines = append(configLines, sc.Text())
	}
	for _, line := range configLines {
		match := submLineRe.FindString(line)
		submodulePath := submPathRe.FindString(match)
		if submodulePath != "" {
			submodulePath = strings.ReplaceAll(submodulePath, `"`, "")
			// format as absolute paths for clarity and safety
			if !filepath.IsAbs(submodulePath) {
				submodulePath = filepath.Join(cfg.Worktree, submodulePath)
			}
			submodulePaths = append(submodulePaths, submodulePath)
		}
	}

	return submodulePaths

}

// Extracts submodule paths from bare repo git config and returns as pointer.
func Submodules() *[]string {
	conf := BareRepoConfig()
	file, err := os.Open(conf)
	if err != nil {
		log.Println("error opening bare repository config:", err)
	}
	defer file.Close()
	submodules := extractSubmodulePaths(file)
	return &submodules

}
