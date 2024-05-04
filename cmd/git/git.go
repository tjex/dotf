package git

import (
	"log"
	"os"
	"os/exec"

	"git.sr.ht/~tjex/dotf/internal/config"
)

func GitCmdExecute(gitArgs []string) {
	argsArray := buildArgsArray(gitArgs)
	cmd := exec.Command("git", argsArray...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

}

func buildArgsArray(gitArgs []string) []string {
	conf := config.UserConfig()
	var argsArray []string

	// force color output
	argsArray = append(argsArray, "-c", "color.status=always")
	argsArray = append(argsArray, conf.RepoFlags...)
	argsArray = append(argsArray, gitArgs...)

	return argsArray
}
