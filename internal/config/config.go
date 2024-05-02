package config

import (
	"bytes"
	toml "github.com/pelletier/go-toml/v2"
	"log"
	"os"
)

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
	conf   Config
)

type Config struct {
	Worktree  string
	GitDir    string
	Origin    string
	Mirror    string
	RepoFlags []string // As a base, targets the bare repo dir and worktree
}

// reads user configuration from a .toml file
func ReadConfig(path string) {
	file, err := os.ReadFile(path)
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

func UserConfig() *Config {
	c := &conf
	return c
}
