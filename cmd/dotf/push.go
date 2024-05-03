package dotf

import (
	"git.sr.ht/~tjex/dotf/cmd/git"
	"git.sr.ht/~tjex/dotf/internal/config"
)

var pushArgsOrigin, pushArgsMirror []string

// push to repositories and mirrors simultaneously
func Push(stdinArgs []string) {
	conf := config.UserConfig()
	c1 := make(chan error, 1)
	c2 := make(chan error, 1)
	pushArgsOrigin = append(pushArgsOrigin, "push", conf.RemoteName)
	pushArgsMirror = append(pushArgsMirror, "push", "--mirror", conf.Mirror)

	go func() {
		// not caring about printing err as git is reporting itself
		err := git.GitCmdRun(pushArgsOrigin)
		c1 <- err
	}()

	go func() {
		err := git.GitCmdRun(pushArgsMirror)
		c2 <- err
	}()
	<-c1
	<-c2

}
