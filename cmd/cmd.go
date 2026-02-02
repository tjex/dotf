package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// A regular exec.Command but stdout and stderr merged and returned as strings.
func Cmd(prog string, args []string) (string, error) {
	var outStd bytes.Buffer
	var outErr bytes.Buffer
	cmd := exec.Command(prog, args...)
	cmd.Stdout = &outStd
	cmd.Stdin = os.Stdin
	cmd.Stderr = &outErr
	cmd.Run()
	if len(outErr.Bytes()) > 0 {
		return "", errors.New(outErr.String())
	}

	return outStd.String(), nil
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
