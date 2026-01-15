package util

import (
	"os"
)

// Get the user set dotf environment. Prioritises shell env variable, DOTF_ENV.
func ResolveEnvironment(env string) string {
	envConfig := env
	envShell := os.Getenv("DOTF_ENV")
	if envShell != "" {
		return envShell
	} else if envConfig != "" {
		return envConfig
	}

	return "default"
}
