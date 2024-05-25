package cmd

import (
	"bytes"
	"log"
	"os"
	"os/exec"

	"git.sr.ht/~tjex/dotf/internal/config"
)

func Cmd(name string, args []string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal("osCmd exited with an error")
	}
}

// A dotf command is a git command with flags set as per the user's
// bare git repository specs

// Execute a regular dotf command (non-concurrent)
func DotfExecute(gitArgs []string) {
	argsArray := buildArgsArray(gitArgs)
	cmd := exec.Command("git", argsArray...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal("DotfCmdExecute exited with an error")
	}

}

// Execute dotf command, but write output to buffer and return
// (for handling via channels). Can only be used where plain output suffices.
func DotfExecuteRoutine(gitArgs []string) string {
	var out bytes.Buffer
	argsArray := buildArgsArray(gitArgs)
	// cmd.Stderr = os.Stderr was writing regular messages to stdout
	// (e.g. printing normal git messages to terminal)?...
	cmd := exec.Command("git", argsArray...)
	cmd.Stdout = &out
	cmd.Stdin = os.Stdin
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		log.Fatal("DotfCmdExecuteRoutine exited with an error")
	}

	return out.String()
}

// Build the arguments array for dotf git call
func buildArgsArray(gitArgs []string) []string {
	conf := config.UserConfig()
	var argsArray []string

	// force color output
	argsArray = append(argsArray, "-c", "color.status=always")
	argsArray = append(argsArray, conf.RepoFlags...)
	argsArray = append(argsArray, gitArgs...)

	return argsArray
}
