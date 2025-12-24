package util

import (
	"os"

	"git.sr.ht/~tjex/dotf/internal/config"
)

var cfg = config.UserConfig()

// Get the user set dotf environment. Prioritises shell env variable, DOTF_ENV.
func ResolveEnvironment() string {
	envConfig := cfg.Env
	envShell := os.Getenv("DOTF_ENV")
	if envShell != "" {
		return envShell
	} else if envConfig != "" {
		return envConfig
	}

	return "default"
}
