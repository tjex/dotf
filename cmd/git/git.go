package git

import (
	"bytes"
	"fmt"
	"git.sr.ht/~tjex/dotf/internal/config"
	"io"
	"log"
	"os/exec"
)

var (
	buf       bytes.Buffer
	logger    = log.New(&buf, "logger: ", log.Lshortfile)
	argsArray []string
)

// prepares bare repo flags for git command and executes
// with provided args
func ExecuteGitCmd(args []string) {
	conf := config.UserConfig()
	argsArray := append(argsArray, conf.RepoFlags...)
	argsArray = append(argsArray, args...)

	cmd := exec.Command("git", argsArray...)

	stderr, err := cmd.StderrPipe()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	slurp, _ := io.ReadAll(stderr)
	out, _ := io.ReadAll(stdout)
	fmt.Printf("%s\n", slurp)
	fmt.Printf("%s\n", out)

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

}
