package dotf

import (
	"fmt"

	"git.sr.ht/~tjex/dotf/cmd"
	"git.sr.ht/~tjex/dotf/internal/config"
)

var pushArgsOrigin, pushArgsMirror []string

// push to repositories and mirrors simultaneously
func Push(stdinArgs []string) {
	conf := config.UserConfig()
	pushArgsOrigin = append(pushArgsOrigin, "push", conf.RemoteName)
	pushArgsMirror = append(pushArgsMirror, "push", "--mirror", conf.Mirror)

	c1 := make(chan string)
	c2 := make(chan string)
	go func() {
		out := cmd.DotfExecuteRoutine(pushArgsOrigin)
		c1 <- out
	}()

	go func() {
		out := cmd.DotfExecuteRoutine(pushArgsMirror)
		c2 <- out
	}()

	out1 := <-c1
	out2 := <-c2

	fmt.Println("---origin---")
	fmt.Println(out1)
	fmt.Println("---mirror---")
	fmt.Println(out2)

}
