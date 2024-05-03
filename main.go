package main

import (
	"bytes"
	"log"
	"os"

	"git.sr.ht/~tjex/dotf/cmd/dotf"
	"git.sr.ht/~tjex/dotf/cmd/git"
	"git.sr.ht/~tjex/dotf/internal/config"
)

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
)

func init() {
	config.ReadConfig("./test/test.toml")
}

func main() {
	stdinArgs := os.Args[1:]

	switch os.Args[1] {
	case "push":
		dotf.Push()
	default:
		// pass all other commands to regular git commands
		git.ExecuteGitCmd(stdinArgs)
	}

}
