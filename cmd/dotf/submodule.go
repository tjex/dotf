package dotf

import (
	"git.sr.ht/~tjex/dotf/internal/config"
)


func Sync() {
	// gitConfig := config.GitConfig()
	config.SubmodulePaths("/Users/tjex/.local/src/dotf/test/test-config")
}

