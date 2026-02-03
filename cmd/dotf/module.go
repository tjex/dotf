package dotf

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"git.sr.ht/~tjex/dotf/cmd"
	"git.sr.ht/~tjex/dotf/internal/config"
	"git.sr.ht/~tjex/dotf/internal/git"
	"git.sr.ht/~tjex/dotf/internal/printer"
	"git.sr.ht/~tjex/dotf/internal/util"
)

var cfg = config.UserConfig()

type ModuleCmd struct {
	Status bool `arg:"--status" default:"false" help:"show 'git status -s' for all modules"`
	Prime  bool `arg:"--prime" default:"false" help:"add and commit all changes to all modules"`
	Push   bool `arg:"--push" default:"false" help:"pushes all modules to their remotes."`
	Pull   bool `arg:"--pull" default:"false" help:"pulls all modules from their remotes."`
	List   bool `arg:"-l,--list" default:"false" help:"list all tracked modules"`
	Edit   bool `arg:"-e, --edit" default:"false" help:"cd into selected module via fzf"`
}

type Modules struct {
	Printer *printer.Printer
	Cmd     *ModuleCmd
}

func errorFmt(repoPath string, err error) error {
	e := fmt.Sprintf("(%s) %v", repoPath, err)
	return errors.New(e)
}

func getModulePaths() []string {
	var paths []string
	env := util.ResolveEnvironment(cfg.Env)

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

func (m *Modules) Run(printer *printer.Printer) error {
	m.Printer = printer
	var err error
	switch {
	case m.Cmd.Status:
		m.status()
	case m.Cmd.Prime:
		err = m.prime()
	case m.Cmd.Push:
		err = m.push()
	case m.Cmd.Pull:
		err = m.pull()
	case m.Cmd.List:
		list()
	case m.Cmd.Edit:
		edit()
	default:
		return errors.New("No flag provided to module command.")
	}
	return err
}

// Add and commit any unstaged changes in all modules.
// There's not real need to check if the repo is dirty. The failure is quick and
// has no side effects.
func (m *Modules) prime() error {
	m.Printer.Println("Checking which modules have uncommitted changes...")
	message := &cfg.BatchCommitMessage

	paths := getModulePaths()

	var wg sync.WaitGroup
	errCh := make(chan error, len(paths))
	wg.Add(len(paths))

	for _, p := range paths {
		go func(p string) {
			defer wg.Done()

			repo, err := util.ExpandPath(p)
			if err != nil {
				errCh <- errorFmt(repo, err)
				return
			}

			report := git.Status(repo)
			// clean repo returns an empty string
			if report != "" {
				m.Printer.Println("\t- Priming", repo)
				git.AddAll(repo)
				git.Commit(repo, message)
			}
		}(p)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		return err
	}
	return nil
}

func (m *Modules) push() error {
	m.Printer.Println("Checking which modules need pushing...")
	paths := getModulePaths()

	var wg sync.WaitGroup
	errCh := make(chan error, len(paths)) // buffered to avoid goroutine blocking

	wg.Add(len(paths))
	for _, p := range paths {
		go func(p string) {
			defer wg.Done()

			repo, err := util.ExpandPath(p)
			if err != nil {
				errCh <- errorFmt(repo, err)
				return
			}

			_, wantsPush, err := git.SyncState(repo)
			if err != nil {
				errCh <- errorFmt(repo, err)
				return
			}

			if wantsPush {
				m.Printer.Println("\t- Pushing", repo)
				git.Push(repo)
			}
		}(p)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		return err
	}

	return nil
}

func (m *Modules) pull() error {
	m.Printer.Println("Checking which modules need pulling...")

	paths := getModulePaths()

	var wg sync.WaitGroup
	errCh := make(chan error, len(paths))
	wg.Add(len(paths))

	for _, p := range paths {
		go func(p string) {
			defer wg.Done()
			repo, err := util.ExpandPath(p)
			if err != nil {
				errCh <- errorFmt(repo, err)
				return
			}
			wantsPull, _, err := git.SyncState(repo)
			if err != nil {
				errCh <- errorFmt(repo, err)
				return
			}
			if wantsPull {
				m.Printer.Println("\t- Pulling:", repo)
				git.Pull(repo)
			}

		}(p)
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		return err
	}

	return nil
}

func (m *Modules) status() error {
	pathsReceived := getModulePaths()

	var wg sync.WaitGroup
	errCh := make(chan error, len(pathsReceived)) // buffered to avoid goroutine blocking

	for _, p := range pathsReceived {
		wg.Add(1)

		go func(p string) {
			defer wg.Done()

			repo, err := util.ExpandPath(p)
			if err != nil {
				errCh <- errorFmt(repo, err)
				return
			}

			report := git.Status(repo)
			if report != "" {
				m.Printer.Println("Changes in", repo+":")
				m.Printer.Println(report)
			}
		}(p)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		return err
	}

	return nil
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
