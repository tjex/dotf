package dotf

import (

	"git.sr.ht/~tjex/dotf/cmd"
)

var pushArgsOrigin, pushArgsMirror []string

// push to repositories and mirrors simultaneously
func Push(remote string) {
	pushArgsOrigin = append(pushArgsOrigin, "push")

	cmd.DotfExecute(pushArgsOrigin)

}
