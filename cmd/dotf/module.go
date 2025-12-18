package dotf

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"git.sr.ht/~tjex/dotf/cmd"
	"git.sr.ht/~tjex/dotf/internal/config"
	"git.sr.ht/~tjex/dotf/internal/git"
	"git.sr.ht/~tjex/dotf/internal/printer"
	"git.sr.ht/~tjex/dotf/internal/util"
)

var cfg = config.UserConfig()

type ModuleCmd struct {
	Prime bool `arg:"--prime" default:"false" help:"add and commit all changes to all modules"`
	Push  bool `arg:"--push" default:"false" help:"pushes all modules to their remotes."`
	Pull  bool `arg:"--pull" default:"false" help:"pulls all modules from their remotes."`
	List  bool `arg:"-l,--list" default:"false" help:"list all tracked modules"`
	Edit  bool `arg:"-e, --edit" default:"false" help:"cd into selected module via fzf"`
}

type Module struct {
	Printer *printer.Printer
	Cmd     *ModuleCmd
}

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

func (m *Module) Run(printer *printer.Printer) error {
	m.Printer = printer
	switch {
	case m.Cmd.Prime:
		m.prime()
	case m.Cmd.Push:
		m.push()
	case m.Cmd.Pull:
		m.pull()
	case m.Cmd.List:
		list()
	case m.Cmd.Edit:
		edit()
	default:
		return errors.New("No flag provided to module command.")
	}
	return nil
}

// Add and commit any unstaged changes in all modules.
// There's not real need to check if the repo is dirty. The failure is quick and
// has no side effects.
func (m *Module) prime() {
	m.Printer.Println("Checking which repos have uncommitted changes...")
	message := &cfg.BatchCommitMessage

	pathsReceived := getModulePaths()
	for _, p := range pathsReceived {

		repo, err := util.ExpandPath(p)
		if err != nil {
			m.Printer.Println(err)
		}

		report := git.Status(repo)
		// clean repo returns an empty string
		if report != "" {
			m.Printer.Println("Priming", repo)
			git.Add(repo)
			git.Commit(repo, message)
		}
	}
}

func (m *Module) push() {
	m.Printer.Println("Checking which repos need pushing...")
	paths := getModulePaths()
	for _, p := range paths {
		repo, err := util.ExpandPath(p)
		if err != nil {
			m.Printer.Println(err)
		}
		wantsPush, _ := git.SyncState(repo)
		if wantsPush {
			m.Printer.Println("Pushing", repo)
			git.Push(repo)
		}
	}
}

func (m *Module) pull() {
	m.Printer.Println("Checking which repos need pulling...")
	paths := getModulePaths()
	for _, p := range paths {
		repo, err := util.ExpandPath(p)
		if err != nil {
			m.Printer.Println(err)
		}
		_, wantsPull := git.SyncState(repo)
		if wantsPull {
			m.Printer.Println("Pulling:", repo)
			git.Pull(repo)
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
func list() {
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

func edit() {
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
