package sync

import (
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
	var err error

	switch {
	case s.Cmd.SyncBare:
		s.syncBare()
	case s.Cmd.SyncModules:
		s.syncModules()
	default:
		err = s.sync()
	}

	return err
}

func (s *Sync) sync() error {
	s.syncBare()
	err := s.syncModules()
	if err != nil {
		return err
	}

	return nil
}

func (s *Sync) syncBare() {
	bare := dotf.Bare{Printer: s.Printer}
	bare.Sync(s.Printer)
}

func (s *Sync) syncModules() error {
	pull := &dotf.ModuleCmd{Pull: true}
	prime := &dotf.ModuleCmd{Prime: true}
	push := &dotf.ModuleCmd{Push: true}
	module := &dotf.Module{Printer: s.Printer}

	module.Cmd = pull
	err := module.Run(s.Printer)
	if err != nil {
		return err
	}

	module.Cmd = prime
	err = module.Run(s.Printer)
	if err != nil {
		return err
	}

	module.Cmd = push
	err = module.Run(s.Printer)
	if err != nil {
		return err
	}

	return nil
}
