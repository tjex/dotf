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
	b.Printer.Println("Syncing bare repository...")

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
		out, err := git.Dotf([]string{"pull"})
		if err != nil {
			return err
		}
		printer.Println(out)
	}

	b.Printer.Println("Priming bare repository...")

	var out string
	out, err = git.Dotf([]string{"add", "-u"}) // add doesnt have a --quiet flag
	if err != nil {
		return err
	}

	printer.Println(out)

	dirty, err := git.UncommittedChanges(bareRepo, worktree)
	if err != nil {
		return err
	}

	if dirty {
		message := &cfg.BatchCommitMessage
		out, err := git.Dotf([]string{"commit", "-m", *message})
		if err != nil {
			return err
		}
		printer.Println(out)
	}

	b.Printer.Println("Pushing bare repository...")
	out, err = git.Dotf([]string{"push"})
	if err != nil {
		return err
	}
	printer.Println(out)
	return nil
}
