package dotf

import (
	"git.sr.ht/~tjex/dotf/internal/git"
	"git.sr.ht/~tjex/dotf/internal/printer"
	"git.sr.ht/~tjex/dotf/internal/util"
)

type Bare struct {
	Printer *printer.Printer
}

func (b *Bare) Sync(printer *printer.Printer) error {
	b.Printer = printer

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
