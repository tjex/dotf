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
	var argsArray []string
	conf := config.UserConfig()
	argsArray = append(argsArray, conf.RepoFlags...)

	switch os.Args[1] {
	case "push":
		dotf.Push()
	default:
		// pass all other commands to regular git commands
		stdinArgs := os.Args[1:]
		argsArray = append(argsArray, stdinArgs...)

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

}
