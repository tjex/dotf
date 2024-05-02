package config

import (
	"bytes"
	"fmt"
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
	Worktree string
	GitDir   string
	Origin   string
	Mirror   string
}

func main() {
	fmt.Println(&buf)
}

// reads a config from a .toml file.
func OpenConfig(path string) Config {
	file, err := os.ReadFile(path)
	if err != nil {
		logger.Print(err)
	}
	conf := parseConfig(file)
	return conf
}

// parse a .toml file for usage.
func parseConfig(config []byte) Config {
	err := toml.Unmarshal(config, &conf)
	if err != nil {
		logger.Print(err)
	}
	return conf
}
