package config

import (
	"bytes"
	"log"
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"
)

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
	conf   Config
)

type Config struct {
	Worktree   string
	GitDir     string
	Origin     string
	Mirror     string
	RepoFlags  []string // As a base, targets the bare repo dir and worktree
	RemoteName string
}

// returns pointer to user config struct
func UserConfig() *Config {
	c := &conf
	return c
}

// reads user configuration from a .toml file
func ReadUserConfig() {
	path := configDir()
	confFile := filepath.Join(path, "config.toml")
	file, err := os.ReadFile(confFile)
	if err != nil {
		logger.Print(err)
	}
	parseConfig(file)

}

// parse a .toml file for usage.
func parseConfig(config []byte) {
	err := toml.Unmarshal(config, &conf)
	if err != nil {
		logger.Print(err)
	}
	buildGitRepoFlags(&conf)
}

// build the bare repo git argument array
func buildGitRepoFlags(conf *Config) {
	var flags []string
	flags = append(flags, conf.GitDir, conf.Worktree)
	conf.RepoFlags = flags
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
