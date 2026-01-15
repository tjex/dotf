package sync

import (
	"git.sr.ht/~tjex/dotf/cmd"
	"git.sr.ht/~tjex/dotf/cmd/dotf"
	"git.sr.ht/~tjex/dotf/internal/printer"
)

type SyncCmd struct {
	Sync        bool `arg:"" help:"sync remote and local with latest changes (i.e, pull all -> push all)."` // default command
	SyncBare    bool `arg:"-b,--bare" default:"false" help:"sync only the bare repository, excluding modules."`
	SyncModules bool `arg:"-m,--modules" default:"false" help:"sync only modules."`
}

type Sync struct {
	Printer *printer.Printer
	Cmd     *SyncCmd
}

func (s *Sync) Run(printer *printer.Printer) error {

	s.Printer = printer

	switch {
	case s.Cmd.SyncBare:
		s.syncBare()
	case s.Cmd.SyncModules:
		s.syncModules()
	default:
		s.sync()
	}

	return nil
}

func (s *Sync) sync() {
	s.syncBare()
	s.syncModules()
}

func (s *Sync) syncBare() {
	cmd.DotfExecute([]string{"pull"}, s.Printer.Quiet)
	cmd.DotfExecute([]string{"push"}, s.Printer.Quiet)
}

func (s *Sync) syncModules() {
	pull := &dotf.ModuleCmd{Pull: true}
	push := &dotf.ModuleCmd{Push: true}
	module := &dotf.Module{Printer: s.Printer}

	module.Cmd = pull
	module.Run(s.Printer)

	module.Cmd = push
	module.Run(s.Printer)
}
