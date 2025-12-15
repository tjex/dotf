package config

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"
)

var (
	buf    bytes.Buffer
	logger = log.New(os.Stdout, "", log.Lshortfile)
	cfg    Config
)

type Config struct {
	Env                string
	Worktree           string
	GitDir             string `toml:"git-dir"`
	Origin             string
	RepoFlags          []string                // Flags for bare repo dir and worktree
	BatchCommitMessage string                  `toml:"batch-commit-message"`
	Modules            map[string]ModuleConfig `toml:"modules"`
}

type ModuleConfig struct {
	Paths []string `toml:"paths"`
}

// returns pointer to user config struct
func UserConfig() *Config {
	return &cfg
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

	err = toml.Unmarshal(file, &cfg)
	if err != nil {
		logger.Print(err)
	}

	// set git flags
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
