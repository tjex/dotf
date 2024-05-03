package dotf

import (
	"git.sr.ht/~tjex/dotf/cmd/git"

	"git.sr.ht/~tjex/dotf/internal/config"
)

var pushArgsOrigin, pushArgsMirror []string

// push to repositories and mirrors simultaneously
func Push(stdinArgs []string) {
	conf := config.UserConfig()
	pushArgsOrigin = append(pushArgsOrigin, "push", conf.RemoteName)
	pushArgsMirror = append(pushArgsMirror, "push", "--mirror", conf.Mirror)

	git.GitCmdRun(pushArgsOrigin)
	git.GitCmdRun(pushArgsMirror)
}
