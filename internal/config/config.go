package config

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"git.sr.ht/~tjex/dotf/internal/util"
	toml "github.com/pelletier/go-toml/v2"
	"github.com/pelletier/go-toml/v2/unstable"
)

var (
	logger = log.New(os.Stdout, "", log.Lshortfile)
	cfg    Config
)

type Config struct {
	Env                string
	Worktree           string
	GitDir             string `toml:"git-dir"`
	Origin             string
	RepoSuffix         string   `toml:"repo-suffix"`
	RepoFlags          []string // Flags for bare repo dir and worktree
	BatchCommitMessage string   `toml:"batch-commit-message"`
	Concurrency        struct {
		BatchSize int `toml:"batch-size"`
	}
	Modules map[string]ModuleConfig `toml:"modules"`
}

type ModuleConfig struct {
	Paths []Module `toml:"paths"`
}

// Module is a tracked directory and, optionally, its remote repository name.
// Repository is used verbatim (suffix is not appended). When empty, the name
// is the directory's last element plus the configured repo-suffix.
type Module struct {
	Path       string `toml:"path"`
	Repo string `toml:"repo"`
}

// UnmarshalTOML accepts either a plain string (the path) or an inline table
// with path and optional repository keys.
func (m *Module) UnmarshalTOML(node *unstable.Node) error {
	switch node.Kind {
	case unstable.String:
		m.Path = string(node.Data)
		return nil
	case unstable.InlineTable:
		it := node.Children()
		for it.Next() {
			kv := it.Node()
			key := kv.Key()
			if !key.Next() {
				continue
			}
			value := string(kv.Value().Data)
			switch string(key.Node().Data) {
			case "path":
				m.Path = value
			case "repo":
				m.Repo = value
			}
		}
		if m.Path == "" {
			return errors.New("module entry has no path")
		}
		return nil
	default:
		return errors.New("module entry must be a string or an inline table")
	}
}

// returns pointer to user config struct
func UserConfig() *Config {
	return &cfg
}

// reads user configuration from a .toml file
func ReadUserConfig() error {
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

	dec := toml.NewDecoder(bytes.NewReader(file))
	dec.EnableUnmarshalerInterface()
	err = dec.Decode(&cfg)
	if err != nil {
		logger.Print(err)
	}

	gitDir, err := util.ExpandPath(cfg.GitDir)
	if err != nil {
		return err
	}
	worktree, err := util.ExpandPath(cfg.Worktree)
	if err != nil {
		return err
	}

	// set git flags
	var flags []string
	flags = append(flags, "--git-dir", gitDir, "--work-tree", worktree)
	cfg.RepoFlags = flags

	setConfigDefaults(&cfg)
	return nil
}

func setConfigDefaults(cfg *Config) {
	// A batch size of 0 will cause concurrency operations to hang.
	if cfg.Concurrency.BatchSize <= 0 {
		cfg.Concurrency.BatchSize = 1000
	}
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
