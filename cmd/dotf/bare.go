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

	wantsPull, _ := git.SyncState(bareRepo)

	if wantsPull {
		out := git.Dotf([]string{"pull"})
		printer.Println(out)
	}

	b.Printer.Println("Priming bare repository...")

	var out string
	out = git.Dotf([]string{"add", "-u"}) // add doesnt have a --quiet flag
	printer.Println(out)

	if dirty := git.UncommittedChanges(bareRepo, worktree); dirty {
		message := &cfg.BatchCommitMessage
		out := git.Dotf([]string{"commit", "-m", *message})
		printer.Println(out)
	}

	b.Printer.Println("Pushing bare repository...")
	out = git.Dotf([]string{"push"})
	printer.Println(out)
	return nil
}
