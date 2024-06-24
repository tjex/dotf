package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"git.sr.ht/~tjex/dotf/internal/config"
)

// A regular exec.Command but stdout and stderr merged and returned as strings.
func Cmd(prog string, args []string) string {
	var out bytes.Buffer
	cmd := exec.Command(prog, args...)
	cmd.Stdout = &out
	cmd.Stdin = os.Stdin
	cmd.Stderr = &out
	cmd.Run() // errors are returned and handled by git itself
	return out.String()
}

// Run fzf by piping arguments.
func CmdFzf(pipeInput []string) (string, error) {
	// deconstruct slice into what fzf sees as a list
	var pipeString string
	for i, val := range pipeInput {
		if i == 0 {
			// otherwise a blank line is entered as the first (0) element
			pipeString = val
			continue
		}
		pipeString = pipeString + "\n" + val
	}

	var result strings.Builder
	fzfPath, err := exec.LookPath("fzf")

	cmd := exec.Command(fzfPath)
	cmd.Stdout = &result
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	_, err = io.WriteString(stdin, pipeString)
	if err != nil {
		return "", err
	}

	err = stdin.Close()
	if err != nil {
		return "", err
	}

	err = cmd.Start()
	if err != nil {
		return "", err
	}

	err = cmd.Wait()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result.String()), nil

}

// open a file or directory with $EDITOR (defaults to vim if no $EDITOR found)
func CmdEditor(path string) {
	var editor, found = os.LookupEnv("EDITOR")
	if !found {
		fmt.Println("No $EDITOR environment variable found, defaulting to vim")
		editor = "vim"
	}
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	var err error
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
	}
	cmd.Run()

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
