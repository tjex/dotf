package dotf

import (
	"fmt"

	"git.sr.ht/~tjex/dotf/cmd"
	"git.sr.ht/~tjex/dotf/internal/printer"
)

type Bare struct {
	Printer *printer.Printer
}

func (b *Bare) Sync(printer *printer.Printer) error {
	b.Printer = printer
	b.Printer.Println("Syncing bare repository...")

	wantsPull, _ := SyncState()

	if wantsPull {
		cmd.DotfExecute([]string{"pull"}, b.Printer.Quiet)
	}

	b.Printer.Println("Priming bare repository...")

	cmd.DotfExecute([]string{"add", "-u"}, false) // add doesnt have a --quiet flag

	if dirty := UncommittedChanges(); dirty {
		message := &cfg.BatchCommitMessage
		cmd.DotfExecute([]string{"commit", "-m", *message}, b.Printer.Quiet)
	}

	b.Printer.Println("Pushing bare repository...")
	cmd.DotfExecute([]string{"push"}, b.Printer.Quiet)
	return nil
}

func SyncState() (bool, bool) {
	cmd.DotfExecute([]string{"fetch", "--quiet"}, true)

	out := cmd.DotfExecute([]string{"rev-list", "--left-right", "--count", "@{upstream}...HEAD"}, true)

	var wantsPull, wantsPush int
	fmt.Sscanf(out, "%d %d", &wantsPull, &wantsPush)

	return wantsPull > 0, wantsPush > 0
}

func UncommittedChanges() bool {
	out := cmd.DotfExecute([]string{"status", "--porcelain"}, false)
	return len(out) > 0
}
