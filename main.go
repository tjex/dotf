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
	config.ReadUserConfig()

}

func main() {
	if len(os.Args) < 2 {
		// Same as simply running "git"
		git.GitCmdExecute(nil)
		return
	}
	stdinArgs := os.Args[1:]

	switch os.Args[1] {
	case "push":
		dotf.Push(stdinArgs)
	case "sync":
		dotf.Sync()
	default:
		git.GitCmdExecute(stdinArgs)
	}

}
