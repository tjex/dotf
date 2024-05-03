package dotf

import (
	"fmt"
	"git.sr.ht/~tjex/dotf/internal/config"
)

// push to repositories and mirrors simultaneously
func Push() {
	conf := config.UserConfig()
	fmt.Println(conf)
	fmt.Println("push from push.go")
}
