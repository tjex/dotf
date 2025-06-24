package dotf

import (
	"fmt"
	"strings"

	"git.sr.ht/~tjex/dotf/cmd"
	"git.sr.ht/~tjex/dotf/internal/config"
	"git.sr.ht/~tjex/dotf/internal/util"
)

var cfg = config.UserConfig()

// Add and commit any unstaged changes in all submodules
func Prime() {
	repos := &cfg.Modules
	message := &cfg.BatchCommitMessage

	for _, s := range *repos {
		// the -C option points to a different cwd for that singular git cmd
		status := []string{"-C", s, "status", "--porcelain"}
		add := []string{"-C", s, "add", "-A"}
		batchCommit := []string{"-C", s, "commit", "-m", *message}
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
	repos := &cfg.Modules
	for _, r := range *repos {
		fmt.Println(r)
	}
}

func Edit() {
	// get submodule paths
	repos := &cfg.Modules

	// return choice from fzf selection
	var choice, err = cmd.CmdFzf(*repos)

	// exit quietly if fzf process is cancelled
	if err != nil && strings.Contains(err.Error(), "exit status 130") {
		return
	}

	// if there's actually an error, print it and exit quietly.
	if err != nil {
		fmt.Println(err)
		return
	}
	choice, err = util.ExpandPath(choice)
	if err != nil {
	    fmt.Println(err)
	}
	cmd.CmdEditor(choice)
}
