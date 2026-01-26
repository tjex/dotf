package dotf

import (
	"git.sr.ht/~tjex/dotf/cmd"
	"git.sr.ht/~tjex/dotf/internal/printer"
)

type Bare struct {
	Printer *printer.Printer
}

func (b *Bare) Sync(printer *printer.Printer) error {
	b.Printer = printer
	b.Printer.Println("Syncing bare repository...")
	cmd.DotfExecute([]string{"pull"}, b.Printer.Quiet)
	err := b.prime()
	if err != nil {
		return err
	}
	cmd.DotfExecute([]string{"push"}, b.Printer.Quiet)
	return nil
}

func (b *Bare) prime() error {
	b.Printer.Println("Priming bare repository...")
	message := &cfg.BatchCommitMessage
	cmd.DotfExecute([]string{"add", "-u"}, b.Printer.Quiet)
	cmd.DotfExecute([]string{"commit", "-m", *message}, b.Printer.Quiet)
	return nil
}
