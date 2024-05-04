package dotf

import (
	"fmt"
	"path/filepath"

	"git.sr.ht/~tjex/dotf/internal/config"
)

func Sync() {
    gitConfig := gitConfig()
	fmt.Println(gitConfig)
}

func gitConfig() string {
	conf := config.UserConfig()
	gitDir := conf.GitDir
	gitConf := filepath.Join(gitDir, "config")
	return gitConf

}
