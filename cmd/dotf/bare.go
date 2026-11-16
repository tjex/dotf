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

func (b *Bare) Sync() error {
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
		cmd := git.Dotf([]string{"pull"})
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	cmd := git.Dotf([]string{"add", "-u"}) // add doesnt have a --quiet flag
	if err := cmd.Run(); err != nil {
		return err
	}

	dirty, err := git.UncommittedChanges(bareRepo, worktree)
	if err != nil {
		return err
	}

	if dirty {
		b.Printer.Println("-> Committing changes.")
		message := &cfg.BatchCommitMessage
		cmd := git.Dotf([]string{"commit", "-m", *message})
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	_, wantsPush, err := git.SyncState(bareRepo)
	if err != nil {
		return err
	}
	if wantsPush {
		b.Printer.Println("-> Pushing.")
		cmd := git.Dotf([]string{"push"})
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bare) Clone() error {
	bareRepo, err := util.ExpandPath(cfg.GitDir)
	if err != nil {
		return err
	}
	cmd := git.Dotf([]string{"clone", "--bare", cfg.Origin, bareRepo})
	if err := cmd.Run(); err != nil {
		return err
	}
	if err := git.Dotf([]string{"config", "--local", "status.showUntrackedFiles", "no"}).Run(); err != nil {
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
	return git.Dotf([]string{"checkout", b.Branch}).Run()
}
