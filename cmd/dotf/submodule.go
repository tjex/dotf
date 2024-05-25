package dotf

import (
	"fmt"

	"git.sr.ht/~tjex/dotf/cmd"
	"git.sr.ht/~tjex/dotf/internal/config"
)

// Add and commit any unstaged changes in all submodules
func CleanAllDirtySubmodules() {
    cfg := config.UserConfig()
    add := []string{"add", "."}
    batchCommit := []string{"commit", "-m", cfg.BatchCommitMessage}
	submodulePaths := config.Submodules()
	for _, s := range submodulePaths {
		status := []string{"-C", s, "status", "--porcelain"}
		out := cmd.Cmd("git", status)
		// clean submodule repos return an empty string
		fmt.Println(batchCommit)
		if out != "" {
			cmd.Cmd("git", add)
			cmd.Cmd("git", batchCommit)
		}
	}
}
