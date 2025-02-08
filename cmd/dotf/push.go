package dotf

import (
	"fmt"

	"git.sr.ht/~tjex/dotf/cmd"
	"git.sr.ht/~tjex/dotf/internal/config"
)

var pushArgsOrigin, pushArgsMirror []string

// push to repositories and mirrors simultaneously
func Push(remote string) {
	conf := config.UserConfig()
	pushArgsOrigin = append(pushArgsOrigin, "push")

	if conf.Mode == "mirror" {
		pushArgsMirror = append(pushArgsMirror, "push", "--mirror", conf.Mirror)
	} else {
		pushArgsMirror = append(pushArgsMirror, "push", conf.Mirror)
	}

	if remote == "" {
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

		fmt.Println(conf.Origin)
		fmt.Println(out1)
		fmt.Println(conf.Mirror)
		fmt.Println(out2)
	} else {
		fmt.Println("execute regular push cmd here")
	}

}
