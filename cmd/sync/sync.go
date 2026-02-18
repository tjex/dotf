package sync

import (
	"errors"

	"git.sr.ht/~tjex/dotf/cmd/dotf"
	"git.sr.ht/~tjex/dotf/internal/printer"
)

type SyncCmd struct {
	SyncAll     bool `arg:"-a,--all" help:"Sync both the bare repository and all modules."`
	SyncBare    bool `arg:"-b,--bare" help:"Sync the bare repository."`
	SyncModules bool `arg:"-m,--modules" help:"Sync all modules."`
}

type Sync struct {
	Printer *printer.Printer
	Cmd     *SyncCmd
}

func (s *Sync) Run(printer *printer.Printer) error {

	s.Printer = printer
	var err error

	switch {
	case s.Cmd.SyncAll:
		err = s.syncAll()
	case s.Cmd.SyncBare:
		err = s.syncBare()
	case s.Cmd.SyncModules:
		err = s.syncModules()
	default:
		return errors.New("No flag passed to sync.")
	}

	return err
}

func (s *Sync) syncAll() error {
	err := s.syncBare()
	if err != nil {
		return err
	}

	err = s.syncModules()
	if err != nil {
		return err
	}

	return nil
}

func (s *Sync) syncBare() error {
	s.Printer.Println("Syncing bare repository...")
	bare := dotf.Bare{Printer: s.Printer}
	err := bare.Sync(s.Printer)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sync) syncModules() error {
	s.Printer.Println("Syncing modules...")
	pull := &dotf.ModuleCmd{Pull: true}
	prime := &dotf.ModuleCmd{Prime: true}
	push := &dotf.ModuleCmd{Push: true}
	module := &dotf.Modules{Printer: s.Printer}

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
