package dotf

import (
	"fmt"
	"strconv"

	"git.sr.ht/~tjex/dotf/cmd"
	"git.sr.ht/~tjex/dotf/internal/config"
)

// Add and commit any unstaged changes in all submodules
func CleanAllDirtySubmodules() {
	submodulePaths := config.Submodules()
	cfg := config.UserConfig()

	message := strconv.Quote(cfg.BatchCommitMessage)
	add := []string{"add", "."}
	batchCommit := []string{"commit", "-m", message}

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
