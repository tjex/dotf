package dotf

import (
	"git.sr.ht/~tjex/dotf/cmd"
	"git.sr.ht/~tjex/dotf/internal/config"
)

func Sync() {
	// gitConfig := config.GitConfig()
	submodulePaths := config.SubmodulePaths("/Users/tjex/.local/src/dotf/test/test-config")

	for _, s := range submodulePaths {
		status := []string{"-C", s, "status"}
		cmd.DotfExecute(status)
	}
}
