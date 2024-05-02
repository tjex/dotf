package main

import (
	"bytes"
	"fmt"
	"git.sr.ht/~tjex/dotf/cmd/dotf"
	"git.sr.ht/~tjex/dotf/internal/config"
	"io"
	"log"
	"os"
	"os/exec"
)

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
)

func init() {
	config.ReadConfig("./test/test.toml")
}

func main() {
	var cmdArgs []string

	arg := os.Args[1]
	switch arg {
	case "push":
		dotf.Push()
	default:
		// pass all other commands to regular git commands
		// following the bare repo user conf entries
		// stdinArgs := os.Args[1:]
		// cmdArgs := append(cmdArgs, gitFlags)
		// cmdArgs = append(cmdArgs, stdinArgs...)

		cmd := exec.Command("git", cmdArgs...)
		stderr, err := cmd.StderrPipe()
		if err != nil {
			fmt.Println(err)
		}

		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}

		slurp, _ := io.ReadAll(stderr)
		fmt.Printf("%s\n", slurp)

		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}
	}

}
