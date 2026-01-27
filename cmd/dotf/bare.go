package dotf

import (
	"git.sr.ht/~tjex/dotf/cmd"
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

	wantsPull, _ := git.SyncState(bareRepo)

	if wantsPull {
		cmd.DotfExecute([]string{"pull"}, b.Printer.Quiet)
	}

	b.Printer.Println("Priming bare repository...")

	cmd.DotfExecute([]string{"add", "-u"}, false) // add doesnt have a --quiet flag

	if dirty := git.UncommittedChanges(bareRepo, worktree); dirty {
		message := &cfg.BatchCommitMessage
		cmd.DotfExecute([]string{"commit", "-m", *message}, b.Printer.Quiet)
	}

	b.Printer.Println("Pushing bare repository...")
	cmd.DotfExecute([]string{"push"}, b.Printer.Quiet)
	return nil
}

