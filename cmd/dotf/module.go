package dotf

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"git.sr.ht/~tjex/dotf/internal/config"
	"git.sr.ht/~tjex/dotf/internal/git"
	"git.sr.ht/~tjex/dotf/internal/printer"
	"git.sr.ht/~tjex/dotf/internal/proc"
	"git.sr.ht/~tjex/dotf/internal/util"
)

var (
	cfg = config.UserConfig()
)

type ModuleCmd struct {
	Status bool `arg:"--status" default:"false" help:"show 'git status -s' for all modules"`
	Prime  bool `arg:"--prime" default:"false" help:"add and commit all changes to all modules"`
	Push   bool `arg:"--push" default:"false" help:"pushes all modules to their remotes."`
	Pull   bool `arg:"--pull" default:"false" help:"pulls all modules from their remotes."`
	List   bool `arg:"-l,--list" default:"false" help:"list all tracked modules"`
	Edit   bool `arg:"-e, --edit" default:"false" help:"cd into selected module via fzf"`
	Clone  bool `arg:"-c, --clone" default:"false" help:"clone modules from default and env groups"`
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

// Run dispatches the flag set on m.Cmd. Other callers should use the
// exported methods directly.
func (m *Modules) Run() error {
	switch {
	case m.Cmd.Status:
		return m.Status()
	case m.Cmd.Prime:
		return m.Prime()
	case m.Cmd.Push:
		return m.Push()
	case m.Cmd.Pull:
		return m.Pull()
	case m.Cmd.List:
		return m.List()
	case m.Cmd.Clone:
		return m.Clone()
	case m.Cmd.Edit:
		return m.Edit()
	default:
		return errors.New("No flag provided to module command.")
	}
}

// Prime adds and commits any unstaged changes in all modules.
// There's not real need to check if the repo is dirty. The failure is quick and
// has no side effects.
func (m *Modules) Prime() error {
	message := &cfg.BatchCommitMessage
	paths := getModulePaths()

	errCh := make(chan error, len(paths))

	var wg sync.WaitGroup

	wg.Add(len(paths))
	for _, p := range paths {
		go func(p string) {
			defer wg.Done()

			repo, err := util.ExpandPath(p)
			if err != nil {
				errCh <- errorFmt(repo, err)
				return
			}

			report, err := git.Status(repo)
			if err != nil {
				errCh <- errorFmt(repo, err)
				return
			}

			// clean repo returns an empty string
			if report != "" {
				m.Printer.Println("-> Priming", repo)
				if _, err := git.AddAll(repo); err != nil {
					errCh <- errorFmt(repo, err)
					return
				}
				if _, err = git.Commit(repo, *message); err != nil {
					errCh <- errorFmt(repo, err)
					return
				}
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

// Push pushes modules to their remotes in parallel.
func (m *Modules) Push() error {
	paths := getModulePaths()
	batchSize := cfg.Concurrency.BatchSize

	errCh := make(chan error, len(paths)) // buffered to avoid goroutine blocking
	sem := make(chan struct{}, batchSize)

	var wg sync.WaitGroup

	for _, p := range paths {
		wg.Add(1)
		sem <- struct{}{}

		go func(p string) {
			defer wg.Done()
			defer func() { <-sem }()

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
				m.Printer.Println("-> Pushing", repo)
				_, err := git.Push(repo)
				if err != nil {
					errCh <- errorFmt(repo, err)
					return
				}
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

// Pull pulls modules from their remotes in parallel.
func (m *Modules) Pull() error {
	paths := getModulePaths()
	batchSize := cfg.Concurrency.BatchSize

	errCh := make(chan error, len(paths))
	sem := make(chan struct{}, batchSize)

	var wg sync.WaitGroup

	for _, p := range paths {
		wg.Add(1)
		sem <- struct{}{}
		go func(p string) {
			defer wg.Done()
			defer func() { <-sem }()
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
				m.Printer.Println("-> Pulling:", repo)
				_, err := git.Pull(repo)
				if err != nil {
					errCh <- errorFmt(repo, err)
					return
				}
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

func (m *Modules) Status() error {
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

			report, err := git.Status(repo)
			if err != nil {
				errCh <- errorFmt(repo, err)
				return
			}

			if report != "" {
				m.Printer.Println("Changes in", repo+":", "\n", report)
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

// List prints the paths of all modules.
func (m *Modules) List() error {
	pathsReceived := getModulePaths()
	paths := append([]string(nil), pathsReceived...)
	sort.Strings(paths)
	for _, path := range paths {
		absPath, err := util.ExpandPath(path)
		if err != nil {
			return err
		}
		fmt.Println(absPath)
	}
	return nil

}

// Edit selects a module via fzf and opens it in $EDITOR.
func (m *Modules) Edit() error {
	paths := getModulePaths()
	sort.Strings(paths)

	choice, err := proc.Fzf(paths)
	if err != nil {
		// fzf exits with 130 when cancelled
		if strings.Contains(err.Error(), "exit status 130") {
			return nil
		}
		return err
	}

	choice, err = util.ExpandPath(choice)
	if err != nil {
		return err
	}
	proc.Editor(choice)
	return nil
}

func (m *Modules) Clone() error {
	paths := getModulePaths()

	var wg sync.WaitGroup
	errCh := make(chan error, len(paths)) // buffered to avoid goroutine blocking

	for _, p := range paths {
		wg.Add(1)

		go func(p string) {
			defer wg.Done()
			repo, err := util.ExpandPath(p)
			if err != nil {
				errCh <- errorFmt(repo, err)
				return
			}
			m.Printer.Println("-> Cloning", repo)
			_, err = git.Clone(repo)
			if err != nil {
				errCh <- errorFmt(repo, err)
				return
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
