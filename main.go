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
	config.ReadConfig("./internal/config/test.toml")
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No CLI arguments passed")
	}
	stdinArgs := os.Args[1:]

	switch os.Args[1] {
	case "push":
		dotf.Push(stdinArgs)
	default:
		git.GitCmdRun(stdinArgs)
	}

}
