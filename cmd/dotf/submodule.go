package dotf

import (
	"git.sr.ht/~tjex/dotf/cmd"
	"git.sr.ht/~tjex/dotf/internal/config"
)

var cfg = config.UserConfig()
var add = []string{"add", "."}
var commit = []string{"commit", "-m", cfg.BatchCommitMessage}

// Add and commit any unstaged changes in all submodules
// "Primes" the submodules for commit and push operations
func SmPrime() {
	submodulePaths := config.SubmodulePaths("/Users/tjex/.local/src/dotf/test/test-config")

	for _, s := range submodulePaths {
		status := []string{"-C", s, "status", "--porcelain"}
		out := cmd.Cmd("git", status)
		// clean submodules repos return an empty string
		if out != "" {
			cmd.Cmd("git", add)
			cmd.Cmd("git", commit)
		}
	}
}
