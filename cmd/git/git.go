package git

import (
	"os"
	"os/exec"

	"git.sr.ht/~tjex/dotf/internal/config"
)

// prepares git command and executes with args passed to function
func GitCmdRun(gitArgs []string) {
	conf := config.UserConfig()
	var argsArray []string

	// force color output
	argsArray = append(argsArray, "-c", "color.status=always")
	argsArray = append(argsArray, conf.RepoFlags...)
	argsArray = append(argsArray, gitArgs...)

	cmd := exec.Command("git", argsArray...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()

}
