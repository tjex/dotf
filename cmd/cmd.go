package cmd

import (
	"bytes"
	"os"
	"os/exec"

	"git.sr.ht/~tjex/dotf/internal/config"
)

// A regular exec.Command but stdout and stderr merged and returned as strings.
func Cmd(name string, args []string) string {
	var out bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &out
	cmd.Stdin = os.Stdin
	cmd.Stderr = &out
	cmd.Run() // errors are returned and handled by git itself
	return out.String()
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
	cmd.Run() // errors are returned and handled by git itself

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
	cmd.Run() // errors are returned and handled by git itself

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
