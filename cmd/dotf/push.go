package dotf

import (
	"sync"

	"git.sr.ht/~tjex/dotf/cmd/git"
	"git.sr.ht/~tjex/dotf/internal/config"
)

var pushArgsOrigin, pushArgsMirror []string

// push to repositories and mirrors simultaneously
func Push(stdinArgs []string) {
	conf := config.UserConfig()
	var wg sync.WaitGroup
	pushArgsOrigin = append(pushArgsOrigin, "push", conf.RemoteName)
	pushArgsMirror = append(pushArgsMirror, "push", "--mirror", conf.Mirror)

	wg.Add(2)
	go func() {
		git.GitCmdExecute(pushArgsOrigin)
		wg.Done()
	}()

	go func() {
		git.GitCmdExecute(pushArgsMirror)
		wg.Done()
	}()

	wg.Wait()

}
