package main

import (
	"bytes"
	"fmt"
	"git.sr.ht/~tjex/dotf/cmd/dotf"
	"git.sr.ht/~tjex/dotf/internal/config"
	"log"
	"os"
	"os/exec"
)

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
)

func main() {
	conf := config.OpenConfig("./test/test.toml")

	logger.Print(conf)
	fmt.Println(&buf)

	arg := os.Args[1]

	switch arg {
	case "push":
		dotf.Push()
	default:
		args := os.Args[1:]
		cmd, err := exec.Command("git", args...).Output()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s", cmd)
	}

}
