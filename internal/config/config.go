package config

import (
	"bufio"
	"bytes"
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
	conf       Config
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
	flags = append(flags, "--git-dir", conf.GitDir, "--work-tree", conf.Worktree)
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

// Retrieve path to bare repo git config
func GitConfig() string {
	conf := UserConfig()
	gitDir := conf.GitDir
	gitConf := filepath.Join(gitDir, "config")
	return gitConf

}

// Extracts submodule paths from bare repo git config
func SubmodulePaths(filepath string) []string {
	var configLines, submodulePaths []string
	file, err := os.Open(filepath)
	if err != nil {
		log.Println("error opening bare repository config:", err)
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

	return submodulePaths

}
