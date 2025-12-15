package util

import (
	"errors"
	"fmt"
	"git.sr.ht/~tjex/dotf/internal/config"
	"os"
)

var cfg = config.UserConfig()

// Get the user set dotf environment. Prioritises shell env variable, DOTF_ENV.
func ResolveEnvironment() string {

	var env string
	envConfig := cfg.Env
	envShell := os.Getenv("DOTF_ENV")
	if envConfig == "" && envShell == "" {
		fmt.Println(errors.New("No dotf environment set."))
		os.Exit(1)
	}
	if envConfig != "" && envShell != "" {
		env = envShell
	} else if envShell != "" {
		env = envShell
	} else {
		env = envConfig
	}

	return env
}
