package main

import (
	"bufio"
	"bytes"
	"fmt"
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
	home, err := os.UserHomeDir()
	if err != nil {
		log.Println("error getting home directory", err)
	}
	confDefaultPath := home + "/config/dotf/config.toml"
	fmt.Println(confDefaultPath)
	// FIXME: can't read conf from XDG_CONFIG_HOME...
	config.ReadConfig("./internal/config/config.toml")

}

func main() {
	if len(os.Args) < 2 {
		// Same as simply running "git"
		git.GitCmdRun(nil)
		return
	}
	stdinArgs := os.Args[1:]
	log.Println(stdinArgs)

	switch os.Args[1] {
	case "push":
		dotf.Push(stdinArgs)
	default:
		out := git.GitCmdRun(stdinArgs)
		f := bufio.NewWriter(os.Stdout)
		defer f.Flush()
		f.Write(out.Bytes())
	}

}
