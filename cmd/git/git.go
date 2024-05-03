package git

import (
	"bytes"
	"log"
	"os"
	"os/exec"

	"git.sr.ht/~tjex/dotf/internal/config"
)

// prepares git command and executes with args passed to function
func GitCmdRun(gitArgs []string) bytes.Buffer {
	conf := config.UserConfig()
	var argsArray []string
	var out bytes.Buffer

	// force color output
	argsArray = append(argsArray, "-c", "color.status=always")
	argsArray = append(argsArray, conf.RepoFlags...)
	argsArray = append(argsArray, gitArgs...)

	cmd := exec.Command("git", argsArray...)
	cmd.Stdout = &out // write stdout to buffer for grouping output between concurrent prints
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	return out

}
