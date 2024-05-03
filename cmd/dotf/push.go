package dotf

import (
	"fmt"

	"git.sr.ht/~tjex/dotf/cmd/git"
	"git.sr.ht/~tjex/dotf/internal/config"
)

var pushArgsOrigin, pushArgsMirror []string

// push to repositories and mirrors simultaneously
func Push(stdinArgs []string) {
	conf := config.UserConfig()
	c1 := make(chan string)
	pushArgsOrigin = append(pushArgsOrigin, "push", conf.RemoteName)
	pushArgsMirror = append(pushArgsMirror, "push", "--mirror", conf.Mirror)

	go func() {
		// not caring about printing err as git is reporting itself
		out := git.GitCmdRun(pushArgsOrigin)
		fmt.Println("repository:", conf.Origin)
		c1 <- out.String()
	}()

	go func() {
		out := git.GitCmdRun(pushArgsMirror)
		c1 <- out.String()
	}()
	<-c1

}
