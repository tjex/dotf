package dotf

import (
	"fmt"
	"os"

	"git.sr.ht/~tjex/dotf/internal/git"
	"git.sr.ht/~tjex/dotf/internal/printer"
	"git.sr.ht/~tjex/dotf/internal/util"
)

type Bare struct {
	Printer *printer.Printer
	Branch  string // branch to check out after cloning; empty keeps the default
}

func (b *Bare) Sync() (err error) {
	defer func() {
		if err != nil {
			err = errorFmt(cfg.GitDir, err)
		}
	}()
	bareRepo, err := util.ExpandPath(cfg.GitDir)
	if err != nil {
		return err
	}

	worktree, err := util.ExpandPath(cfg.Worktree)
	if err != nil {
		return err
	}

	wantsPull, _, err := git.SyncState(bareRepo)
	if err != nil {
		return err
	}

	if wantsPull {
		if err := git.RunDotf([]string{"pull"}); err != nil {
			return err
		}
	}

	if err := git.RunDotf([]string{"add", "-u"}); err != nil {
		return err
	}

	dirty, err := git.UncommittedChanges(bareRepo, worktree)
	if err != nil {
		return err
	}

	if dirty {
		b.Printer.Println("-> Committing changes.")
		message := &cfg.BatchCommitMessage
		if err := git.RunDotf([]string{"commit", "-m", *message}); err != nil {
			return err
		}
	}

	_, wantsPush, err := git.SyncState(bareRepo)
	if err != nil {
		return err
	}
	if wantsPush {
		b.Printer.Println("-> Pushing.")
		if err := git.RunDotf([]string{"push"}); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bare) Clone() (err error) {
	defer func() {
		if err != nil {
			err = errorFmt(cfg.GitDir, err)
		}
	}()
	bareRepo, err := util.ExpandPath(cfg.GitDir)
	if err != nil {
		return err
	}
	if err := git.RunDotf([]string{"clone", "--bare", cfg.Origin, bareRepo}); err != nil {
		return err
	}
	if err := git.RunDotf([]string{"config", "--local", "status.showUntrackedFiles", "no"}); err != nil {
		return err
	}
	if b.Branch == "" {
		return nil
	}
	if !git.BareBranchExists(bareRepo, b.Branch) {
		fmt.Fprintf(os.Stderr, "Branch %s does not exist in repo %s. Cloned the default branch.\n", b.Branch, bareRepo)
		return nil
	}
	b.Printer.Println("-> Checking out", b.Branch, "in", bareRepo)
	return git.RunDotf([]string{"checkout", b.Branch})
}
