package dotf

import (
	"fmt"
	"sort"
	"strings"

	"git.sr.ht/~tjex/dotf/cmd"
	"git.sr.ht/~tjex/dotf/internal/config"
	"git.sr.ht/~tjex/dotf/internal/git"
	"git.sr.ht/~tjex/dotf/internal/util"
)

var cfg = config.UserConfig()

func getModulePaths() []string {
	var paths []string
	env := util.ResolveEnvironment()

	modules := &cfg.Modules
	for name, m := range *modules {
		// modules.default are always included
		if name != env && name != "default" {
			continue
		}
		for _, p := range m.Paths {
			paths = append(paths, p)
		}
	}

	return paths
}

// Add and commit any unstaged changes in all modules.
// There's not real need to check if the repo is dirty. The failure is quick and
// has no side effects.
func Prime() {
	message := &cfg.BatchCommitMessage

	pathsReceived := getModulePaths()
	for _, p := range pathsReceived {

		repo, err := util.ExpandPath(p)
		if err != nil {
			fmt.Println(err)
		}

		report := git.Status(repo)
		// clean repo returns an empty string
		if report != "" {
			fmt.Println("Priming", repo)
			git.Add(repo)
			git.Commit(repo, message)
		}
	}
}

func Push() {
	paths := getModulePaths()
	for _, p := range paths {
		repo, err := util.ExpandPath(p)
		if err != nil {
			fmt.Println(err)
		}
		wantsPush, _ := git.SyncState(repo)
		if wantsPush {
			fmt.Println("Pushing", repo)
			git.Push(repo)
		} else {
			fmt.Println("Checking", repo)
		}
	}
}

func Pull() {
	paths := getModulePaths()
	for _, p := range paths {
		repo, err := util.ExpandPath(p)
		if err != nil {
			fmt.Println(err)
		}
		_, wantsPull := git.SyncState(repo)
		if wantsPull {
			fmt.Println("Pulling", repo)
			git.Pull(repo)
		} else {
			fmt.Println("Checking", repo)
		}
	}
}

// TODO
func Sync() {
	// pull
	// prime
	// push
	// needs to handle merge conflict reports etc
}

// Return paths to all submodules
func List() {
	pathsReceived := getModulePaths()
	var paths []string
	for _, path := range pathsReceived {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Println(path)
	}

}

func Edit() {
	// get submodule paths
	var paths []string
	pathsReceived := getModulePaths()
	for _, p := range pathsReceived {
		paths = append(paths, p)
	}

	sort.Strings(paths)
	// return choice from fzf selection
	var choice, err = cmd.CmdFzf(paths)

	// exit quietly if fzf process is cancelled
	if err != nil && strings.Contains(err.Error(), "exit status 130") {
		return
	}

	// if there's actually an error, print it and exit quietly.
	if err != nil {
		fmt.Println(err)
		return
	}
	choice, err = util.ExpandPath(choice)
	if err != nil {
		fmt.Println(err)
	}
	cmd.CmdEditor(choice)
}
