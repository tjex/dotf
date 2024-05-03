package git

import (
	"fmt"
	"git.sr.ht/~tjex/dotf/internal/config"
	"io"
	"log"
	"os/exec"
)

// prepares git command and executes with args passed to function
func ExecuteGitCmd(args []string) {
	conf := config.UserConfig()
	var argsArray []string

	// force color output
	argsArray = append(argsArray, "-c", "color.status=always")
	argsArray = append(argsArray, conf.RepoFlags...)
	argsArray = append(argsArray, args...)

	cmd := exec.Command("git", argsArray...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	slurpErr, _ := io.ReadAll(stderr)
	slurpOut, _ := io.ReadAll(stdout)

	fmt.Printf("%s\n", slurpErr)
	fmt.Printf("%s\n", slurpOut)

	// be sure cmd finishes
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

}
