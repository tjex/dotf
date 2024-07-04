package dotf

import (
	"fmt"
	"strings"

	"git.sr.ht/~tjex/dotf/cmd"
	"git.sr.ht/~tjex/dotf/internal/config"
)

// Add and commit any unstaged changes in all submodules
func Prime() {
	submodulePaths := config.Submodules()
	cfg := config.UserConfig()

	message := cfg.BatchCommitMessage

	for _, s := range *submodulePaths {
		status := []string{"-C", s, "status", "--porcelain"}
		add := []string{"-C", s, "add", "-A"}
		batchCommit := []string{"-C", s, "commit", "-m", message}
		report := cmd.Cmd("git", status)
		// clean submodule repos return an empty string
		if report != "" {
			cmd.Cmd("git", add)
			cmd.Cmd("git", batchCommit)
		}
	}
}

// Return paths to all submodules
func List() {
	submodulePaths := config.Submodules()
	for _, submodule := range *submodulePaths {
		fmt.Println(submodule)
	}
}

func Edit() {
	// get submodule paths
	submodulePaths := config.Submodules()
	var pathsDeref = *submodulePaths

	// return choice from fzf selection
	var choice, err = cmd.CmdFzf(pathsDeref)

	// exit quietly if fzf process is cancelled
	if err != nil && strings.Contains(err.Error(), "exit status 130") {
		return
	}

	// if there's actually an error, print it and exit quietly.
	if err != nil {
		fmt.Println(err)
		return
	}
	cmd.CmdEditor(choice)
}
